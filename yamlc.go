package yamlc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

// CommentStyle 注释风格枚举
type CommentStyle int

const (
	// StyleTop 注释放在字段头顶（标准YAML风格）
	StyleTop CommentStyle = iota
	// StyleInline 注释放在字段后面（行内注释，需要对齐）
	StyleInline
	// StyleSmart 智能注释：如果没有子元素就加到末尾(需要对齐)，有的就加在头顶
	StyleSmart
	// StyleCompact 紧凑风格：注释和值在同一行，单空格分隔
	StyleCompact
	// StyleMinimal 最小风格：只显示字段和值，无注释
	StyleMinimal
	// StyleVerbose 详细风格：显示字段名(类型):详细注释
	StyleVerbose
	// StyleSpaced 间隔风格：头顶风格的延续，在字段之间添加空行分隔
	StyleSpaced
	// StyleGrouped 分组风格：常规数据类型就放在一起，当同级遇到复杂数据类型，组间有空行分隔
	StyleGrouped
	// StyleSectioned 分节风格：同级的标签常规数据类型聚合在第一个标签顶部，使用多注释
	StyleSectioned
	// StyleDoc 文档风格：使用分隔线和标题，详细的类型信息
	StyleDoc
	// StyleSeparate 分离风格：注释和值分离，所有注释在前面，值在后面
	StyleSeparate
)

// GlobalCommentStyle 全局注释风格设置
var GlobalCommentStyle = StyleTop

func init() {
	GlobalCommentStyle = StyleTop
}

func GetStyle() CommentStyle {
	return GlobalCommentStyle
}

func GetAllStyle() []CommentStyle {
	return []CommentStyle{
		StyleTop,
		StyleInline,
		StyleSmart,
		StyleCompact,
		StyleMinimal,
		StyleVerbose,
		StyleSpaced,
		StyleGrouped,
		StyleSectioned,
		StyleDoc,
		StyleSeparate,
	}
}

func GetStyleString(style int) string {
	switch style {
	case 0:
		return "top"
	case 1:
		return "inline"
	case 2:
		return "smart"
	case 3:
		return "compact"
	case 4:
		return "minimal"
	case 5:
		return "verbose"
	case 6:
		return "spaced"
	case 7:
		return "grouped"
	case 8:
		return "sectioned"
	case 9:
		return "doc"
	case 10:
		return "separate"
	}
	return "smart"
}

type Option func(*Options)

type Options struct {
	Style    CommentStyle
	Comments []map[string]string
}

func WithStyle(style CommentStyle) Option {
	return func(o *Options) {
		o.Style = style
	}
}

func WithComment(comments map[string]string) Option {
	return func(o *Options) {
		o.Comments = append(o.Comments, comments)
	}
}

// FieldInfo 字段信息结构
type FieldInfo struct {
	Name        string
	Comment     string
	Field       reflect.Value
	FieldType   reflect.StructField
	HasChildren bool
	FieldPath   string
}

// Gen 生成YAML内容
func Gen(v interface{}, opts ...Option) ([]byte, error) {
	options := &Options{
		Style:    GlobalCommentStyle,
		Comments: make([]map[string]string, 0),
	}

	for _, opt := range opts {
		opt(options)
	}

	if v == nil {
		return nil, fmt.Errorf("input value cannot be nil")
	}

	var result []byte
	if options.Style == StyleMinimal {
		yamlData, err := generateMinimalStyleField(v)
		if err != nil {
			return nil, fmt.Errorf("failed to generate YAML content: %w", err)
		}
		result = []byte(yamlData)
	} else {

		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return nil, fmt.Errorf("input pointer cannot be nil")
			}
			val = val.Elem()
		}

		var buf bytes.Buffer

		content, err := generateValue(val, "", 0, options)
		if err != nil {
			return nil, fmt.Errorf("failed to generate YAML content: %w", err)
		}

		buf.WriteString(content)

		result = buf.Bytes()
	}
	// 严格的YAML格式验证
	if err := ValidateYAML(result); err != nil {
		return nil, fmt.Errorf("generated YAML validation failed: %w", err)
	}

	return result, nil
}

// Write 写入到io.Writer
func Write(w io.Writer, v interface{}, opts ...Option) error {
	if w == nil {
		return fmt.Errorf("writer cannot be nil")
	}

	data, err := Gen(v, opts...)
	if err != nil {
		return err
	}

	n, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("incomplete write: wrote %d bytes, expected %d", n, len(data))
	}

	return nil
}

// WriteFile 写入到文件
func WriteFile(filename string, v interface{}, opts ...Option) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", filename, err)
	}
	defer file.Close()

	return Write(file, v, opts...)
}

// ValidateYAML 使用yaml.v3进行严格的YAML格式验证
func ValidateYAML(data []byte) error {
	var result interface{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(false) // 允许未知字段

	err := decoder.Decode(&result)
	if err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}

	return nil
}

// ValidateStructure 验证YAML结构的完整性
func ValidateStructure(data []byte) error {
	lines := strings.Split(string(data), "\n")
	var indentStack []int

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// 跳过空行和注释
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// 计算缩进
		indent := len(line) - len(strings.TrimLeft(line, " "))

		// 验证缩进是否为偶数（YAML标准）
		if indent%2 != 0 {
			return fmt.Errorf("invalid indentation at line %d: indentation must be even number of spaces", lineNum)
		}

		// 验证缩进层级的一致性
		if len(indentStack) == 0 {
			indentStack = append(indentStack, indent)
		} else {
			currentLevel := indent / 2
			if currentLevel > len(indentStack) {
				return fmt.Errorf("invalid indentation jump at line %d: too many levels", lineNum)
			}
			// 调整堆栈
			indentStack = indentStack[:currentLevel+1]
			if currentLevel < len(indentStack) {
				indentStack[currentLevel] = indent
			} else {
				indentStack = append(indentStack, indent)
			}
		}

		// 验证键值对格式
		if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "-") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid key-value format at line %d: %s", lineNum, trimmed)
			}

			key := strings.TrimSpace(parts[0])
			if key == "" {
				return fmt.Errorf("empty key at line %d", lineNum)
			}

			// 验证键名格式
			if err := validateKeyName(key); err != nil {
				return fmt.Errorf("invalid key name at line %d: %w", lineNum, err)
			}
		}
	}

	return nil
}

// validateKeyName 验证键名格式
func validateKeyName(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// 如果是引用字符串，直接通过
	if (strings.HasPrefix(key, `"`) && strings.HasSuffix(key, `"`)) ||
		(strings.HasPrefix(key, `'`) && strings.HasSuffix(key, `'`)) {
		return nil
	}

	// 检查是否包含YAML特殊字符
	if strings.ContainsAny(key, ":#{}[]&*!|>\"'%?@`") {
		return fmt.Errorf("key contains YAML special characters: %s", key)
	}

	return nil
}

// Validate 旧的验证函数，保持向后兼容
func Validate(v []byte) error {
	return ValidateYAML(v)
}

// generateValue 递归生成YAML值
func generateValue(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	if !val.IsValid() {
		return "null", nil
	}

	switch val.Kind() {
	case reflect.Struct:
		return generateStruct(val, fieldPath, indent, options)
	case reflect.Map:
		return generateMap(val, fieldPath, indent, options)
	case reflect.Slice, reflect.Array:
		return generateSlice(val, fieldPath, indent, options)
	case reflect.String:
		return generateString(val, fieldPath, indent, options)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return generateInt(val, fieldPath, indent, options)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return generateUint(val, fieldPath, indent, options)
	case reflect.Float32, reflect.Float64:
		return generateFloat(val, fieldPath, indent, options)
	case reflect.Bool:
		return generateBool(val, fieldPath, indent, options)
	case reflect.Ptr:
		if val.IsNil() {
			return "null", nil
		}
		return generateValue(val.Elem(), fieldPath, indent, options)
	case reflect.Interface:
		if val.IsNil() {
			return "null", nil
		}
		return generateValue(val.Elem(), fieldPath, indent, options)
	default:
		if val.CanInterface() {
			return fmt.Sprintf("%v", val.Interface()), nil
		}
		return "null", nil
	}
}

// generateStruct 生成结构体YAML
func generateStruct(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	typ := val.Type()
	fields := collectFieldInfo(val, typ, fieldPath, options)

	if len(fields) == 0 {
		return " {}\n", nil
	}

	var result string
	var err error

	switch options.Style {
	case StyleDoc:
		result, err = generateStructDoc(fields, indent, options)
	case StyleSeparate:
		result, err = generateStructSeparate(fields, indent, options)
	case StyleSectioned:
		result, err = generateStructSectioned(fields, indent, options)
	default:
		result, err = generateStructDefault(fields, indent, options)
	}

	if err != nil {
		return "", err
	}

	result = result + "\n"

	return result, nil
}

// collectFieldInfo 收集字段信息
func collectFieldInfo(val reflect.Value, typ reflect.Type, fieldPath string, options *Options) []FieldInfo {
	var fields []FieldInfo

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		fieldName := getFieldName(fieldType)
		if fieldName == "-" {
			continue
		}

		currentFieldPath := buildFieldPath(fieldPath, fieldName)
		comment := getComment(fieldType, currentFieldPath, options)
		hasChildren := hasChildren(field)

		fields = append(fields, FieldInfo{
			Name:        fieldName,
			Comment:     comment,
			Field:       field,
			FieldType:   fieldType,
			HasChildren: hasChildren,
			FieldPath:   currentFieldPath,
		})
	}

	return fields
}

// generateStructDoc 生成文档风格的结构体
func generateStructDoc(fields []FieldInfo, indent int, options *Options) (string, error) {
	var result strings.Builder
	indentStr := strings.Repeat("  ", indent)

	// 生成文档头部注释块

	result.WriteString(fmt.Sprintf("%s############################################\n", indentStr))
	for _, field := range fields {
		if field.Comment != "" {
			typeStr := field.Field.Type().String()
			result.WriteString(fmt.Sprintf("%s# %s(%s):%s\n", indentStr, field.Name, typeStr, field.Comment))
		}
		if field.HasChildren {
			break
		}
	}
	result.WriteString(fmt.Sprintf("%s###########################################\n\n", indentStr))

	// 生成字段
	for i, field := range fields {

		result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))

		if field.HasChildren {
			result.WriteString("\n")
			fieldValue, err := generateValue(field.Field, field.FieldPath, indent+1, options)
			if err != nil {
				return "", err
			}
			result.WriteString(fieldValue)
		} else {
			fieldValue, err := generateValue(field.Field, field.FieldPath, indent+1, options)
			if err != nil {
				return "", err
			}
			result.WriteString(fieldValue)
		}

		if i < len(fields)-1 {
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

// generateStructSeparate 生成分离风格的结构体
func generateStructSeparate(fields []FieldInfo, indent int, options *Options) (string, error) {
	var result strings.Builder

	// 如果是顶层，先生成所有注释
	if indent == 0 {
		result.WriteString("############################################\n")
		generateAllComments(&result, fields, 0, "")
		result.WriteString("###########################################\n\n")
	}

	// 生成字段值
	indentStr := strings.Repeat("  ", indent)
	for i, field := range fields {
		result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))

		if field.HasChildren {
			result.WriteString("\n")
			fieldValue, err := generateValue(field.Field, field.FieldPath, indent+1, options)
			if err != nil {
				return "", err
			}
			result.WriteString(fieldValue)
		} else {
			fieldValue, err := generateValue(field.Field, field.FieldPath, indent+1, options)
			if err != nil {
				return "", err
			}
			result.WriteString(fieldValue)
		}

		if i < len(fields)-1 {
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

// generateAllComments 递归生成所有注释
func generateAllComments(result *strings.Builder, fields []FieldInfo, indent int, prefix string) {
	// fmt.Println("generateAllComments", fields)
	for _, field := range fields {
		// if field.Comment != "" {
		typeStr := field.Field.Type().String()
		indentStr := strings.Repeat("  ", indent)
		result.WriteString(fmt.Sprintf("# %s%s(%s):%s\n", indentStr, field.Name, typeStr, field.Comment))
		// }

		// 如果有子结构，递归生成子注释
		if field.HasChildren {
			var subFields []FieldInfo

			switch field.Field.Kind() {
			case reflect.Struct:
				// 结构体类型，直接收集字段信息
				subFields = collectFieldInfo(field.Field, field.Field.Type(), field.FieldPath, &Options{})
			case reflect.Ptr:
				// 指针类型，解引用后收集字段信息
				if !field.Field.IsNil() {
					elem := field.Field.Elem()
					if elem.Kind() == reflect.Struct {
						subFields = collectFieldInfo(elem, elem.Type(), field.FieldPath, &Options{})
					}
				}
			case reflect.Slice, reflect.Array:
				// 切片/数组类型，检查第一个元素
				if field.Field.Len() > 0 {
					firstItem := field.Field.Index(0)
					if firstItem.Kind() == reflect.Struct {
						subFields = collectFieldInfo(firstItem, firstItem.Type(), field.FieldPath+"[0]", &Options{})
					} else if firstItem.Kind() == reflect.Ptr && !firstItem.IsNil() {
						elem := firstItem.Elem()
						if elem.Kind() == reflect.Struct {
							subFields = collectFieldInfo(elem, elem.Type(), field.FieldPath+"[0]", &Options{})
						}
					}
				}
			case reflect.Map:
				// 映射类型，检查第一个值
				if field.Field.Len() > 0 {
					iter := field.Field.MapRange()
					if iter.Next() {
						value := iter.Value()
						if value.Kind() == reflect.Struct {
							subFields = collectFieldInfo(value, value.Type(), field.FieldPath+"[key]", &Options{})
						} else if value.Kind() == reflect.Ptr && !value.IsNil() {
							elem := value.Elem()
							if elem.Kind() == reflect.Struct {
								subFields = collectFieldInfo(elem, elem.Type(), field.FieldPath+"[key]", &Options{})
							}
						}
					}
				}
			}

			if len(subFields) > 0 {
				generateAllComments(result, subFields, indent+1, prefix+"    ")
			}
		}
	}
}

// generateStructSectioned 生成分节风格的结构体
func generateStructSectioned(fields []FieldInfo, indent int, options *Options) (string, error) {
	var result strings.Builder
	indentStr := strings.Repeat("  ", indent)

	type FieldInfoArr struct {
		Fields   []FieldInfo
		isSimple bool
	}

	// 分组字段：简单字段和复杂字段
	old_hc_label := true
	var fieldInfoArrs []FieldInfoArr
	for _, field := range fields {
		if field.HasChildren {
			fieldInfoArrs = append(fieldInfoArrs, FieldInfoArr{Fields: []FieldInfo{field}, isSimple: false})
		} else {
			if old_hc_label {
				fieldInfoArrs = append(fieldInfoArrs, FieldInfoArr{Fields: []FieldInfo{field}, isSimple: true})
				old_hc_label = false
			} else {
				fieldInfoArrs[len(fieldInfoArrs)-1].Fields = append(fieldInfoArrs[len(fieldInfoArrs)-1].Fields, field)
			}
		}
	}

	for _, fieldInfoArr := range fieldInfoArrs {

		// 先处理简单字段，在第一个字段上方集中显示注释
		if fieldInfoArr.isSimple {
			for _, field := range fieldInfoArr.Fields {
				if field.Comment != "" {
					result.WriteString(fmt.Sprintf("%s# %s\n", indentStr, field.Comment))
				}
			}
			result.WriteString("\n")

			for i, field := range fieldInfoArr.Fields {
				result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))
				fieldValue, err := generateValue(field.Field, field.FieldPath, indent+1, options)
				if err != nil {
					return "", err
				}
				result.WriteString(" " + strings.TrimSpace(fieldValue))
				if i < len(fieldInfoArr.Fields)-1 {
					result.WriteString("\n")
				}
			}
			if len(fieldInfoArr.Fields) > 0 {
				result.WriteString("\n")
			}

		}
		if !fieldInfoArr.isSimple {

			// 再处理复杂字段
			result.WriteString("\n")
			for i, field := range fieldInfoArr.Fields {
				if field.Comment != "" {
					result.WriteString(fmt.Sprintf("%s# %s\n", indentStr, field.Comment))
				}
				result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))

				fieldValue, err := generateValue(field.Field, field.FieldPath, indent+1, options)
				if err != nil {
					return "", err
				}
				//如果fieldValue第一行是注释就换行
				if strings.HasPrefix(strings.TrimSpace(fieldValue), "#") {
					result.WriteString("\n")
				}
				result.WriteString(fieldValue)

				if i < len(fieldInfoArr.Fields)-1 {
					result.WriteString("\n")
				}
			}
		}
	}

	return result.String(), nil
}

// generateStructDefault 生成默认风格的结构体
func generateStructDefault(fields []FieldInfo, indent int, options *Options) (string, error) {
	var result strings.Builder
	maxFieldNameLen := calculateMaxFieldNameLen(fields)

	for i, field := range fields {
		if err := generateFieldWithComment(&result, field, indent, options.Style, maxFieldNameLen, options); err != nil {
			return "", err
		}

		// 添加字段间间隔
		if shouldAddSpacing(options.Style, i, len(fields)) {
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

// calculateMaxFieldNameLen 计算最大字段名长度
func calculateMaxFieldNameLen(fields []FieldInfo) int {
	maxLen := 0
	for _, field := range fields {
		// 只计算字段名长度，不包括字段值，因为字段值可能很长
		fieldNameLen := len(field.Name) + 1 // +1 for colon
		if fieldNameLen > maxLen {
			maxLen = fieldNameLen
		}
	}
	return maxLen
}

// shouldAddSpacing 判断是否应该添加间距
func shouldAddSpacing(style CommentStyle, currentIndex, totalFields int) bool {
	if currentIndex >= totalFields-1 {
		return false
	}

	return style == StyleSpaced || style == StyleGrouped || style == StyleSectioned
}

// generateFieldWithComment 生成带注释的字段
func generateFieldWithComment(result *strings.Builder, field FieldInfo, indent int,
	commentStyle CommentStyle, maxFieldNameLen int, options *Options) error {

	indentStr := strings.Repeat("  ", indent)

	// 智能风格的动态调整
	if commentStyle == StyleSmart {
		if field.HasChildren {
			commentStyle = StyleTop
		} else {
			commentStyle = StyleInline
		}
	}

	switch commentStyle {
	case StyleTop:
		return generateTopStyleField(result, field, indentStr, options)
	case StyleInline:
		return generateInlineStyleField(result, field, indentStr, maxFieldNameLen, options)
	case StyleCompact:
		return generateCompactStyleField(result, field, indentStr, options)
	case StyleVerbose:
		return generateVerboseStyleField(result, field, indentStr, options)
	case StyleSpaced, StyleGrouped:
		return generateTopStyleField(result, field, indentStr, options)
	default:
		return generateTopStyleField(result, field, indentStr, options)
	}
}

// generateTopStyleField 生成顶部风格字段
func generateTopStyleField(result *strings.Builder, field FieldInfo, indentStr string, options *Options) error {
	if field.Comment != "" {
		result.WriteString(fmt.Sprintf("%s# %s\n", indentStr, field.Comment))
	}
	result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))

	return generateFieldValue(result, field, indentStr, options)
}

// generateInlineStyleField 生成内联风格字段
func generateInlineStyleField(result *strings.Builder, field FieldInfo, indentStr string, maxFieldNameLen int, options *Options) error {
	maxFieldNameLen = maxFieldNameLen + 30
	fieldNamePart := field.Name + ":"
	currentFieldNameLen := getDisplayWidth(fieldNamePart)

	// 处理复杂类型（有子字段的情况）
	if field.HasChildren {
		if field.Comment != "" {
			// 检查是否为空容器
			if isEmpty := isEmptyContainer(field.Field); isEmpty {
				emptyValue := getEmptyContainerValue(field.Field)
				alignSpaces := maxFieldNameLen - currentFieldNameLen - getDisplayWidth(emptyValue) + 2
				if alignSpaces < 1 {
					alignSpaces = 1
				}
				result.WriteString(fmt.Sprintf("%s%s %s%s# %s\n",
					indentStr, fieldNamePart, emptyValue,
					strings.Repeat(" ", alignSpaces), field.Comment))
				return nil
			}

			// 非空复杂类型，注释放在同一行
			alignSpaces := maxFieldNameLen - currentFieldNameLen + 2
			if alignSpaces < 1 {
				alignSpaces = 1
			}
			result.WriteString(fmt.Sprintf("%s%s%s# %s",
				indentStr, fieldNamePart,
				strings.Repeat(" ", alignSpaces), field.Comment))
		} else {
			result.WriteString(fmt.Sprintf("%s%s ", indentStr, fieldNamePart))
		}
		return generateFieldValue(result, field, indentStr, options)
	}

	indent := 0
	if field.Field.Kind() == reflect.Slice || field.Field.Kind() == reflect.Array {
		hasVisibleChildren := field.HasChildren || (field.Field.Kind() == reflect.Slice && field.Field.Len() > 0) ||
			(field.Field.Kind() == reflect.Array && field.Field.Len() > 0)
		if hasVisibleChildren {
			// 复杂类型使用顶部注释
			if field.Comment != "" {
				fieldNameAndValueWidth := getDisplayWidth(indentStr + field.Name + ": ")
				alignSpaces := maxFieldNameLen + getDisplayWidth(indentStr) - fieldNameAndValueWidth
				result.WriteString(fmt.Sprintf("%s%s:%s# %s", indentStr, field.Name, strings.Repeat(" ", alignSpaces), field.Comment))
			} else {
				result.WriteString(fmt.Sprintf("%s%s:", indentStr, field.Name))
			}
			indent = 1
		} else {
			result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))
		}
		// 生成字段值
		fieldValue, err := generateValue(field.Field, field.FieldPath, indent, options)
		if err != nil {
			return err
		}
		fieldValue = strings.TrimRight(fieldValue, "\n")

		if field.Comment != "" {
			if hasVisibleChildren {
				result.WriteString(fmt.Sprintf("%s\n", fieldValue))
			} else {
				// 计算对齐空格 - 使用实际的字段名和值长度
				fieldNameAndValueWidth := getDisplayWidth(indentStr + field.Name + ": " + fieldValue)
				alignSpaces := maxFieldNameLen + getDisplayWidth(indentStr) - fieldNameAndValueWidth + 2
				if alignSpaces < 1 {
					alignSpaces = 1
				}
				result.WriteString(fmt.Sprintf("%s%s# %s\n", fieldValue, strings.Repeat(" ", alignSpaces), field.Comment))
			}
		} else {
			result.WriteString(fmt.Sprintf("%s\n", fieldValue))
		}
		return nil
	} else {
		result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))
	}

	// 生成字段值
	fieldValue, err := generateValue(field.Field, field.FieldPath, 0, options)
	if err != nil {
		return err
	}
	fieldValue = strings.TrimRight(fieldValue, "\n")

	// 计算注释对齐
	if field.Comment != "" {
		// 计算实际的字段行宽度
		actualFieldLine := fmt.Sprintf("%s%s: %s", indentStr, field.Name, fieldValue)
		fieldLineWidth := getDisplayWidth(actualFieldLine)

		// 计算对齐空格
		targetWidth := getDisplayWidth(indentStr) + maxFieldNameLen
		alignSpaces := targetWidth - fieldLineWidth + 2
		if alignSpaces < 1 {
			alignSpaces = 1
		}

		result.WriteString(fmt.Sprintf("%s%s# %s\n",
			fieldValue, strings.Repeat(" ", alignSpaces), field.Comment))
	} else {
		result.WriteString(fmt.Sprintf("%s\n", fieldValue))
	}

	return nil
}

// isEmptyContainer 检查容器是否为空
func isEmptyContainer(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return field.Len() == 0
	case reflect.Struct:
		// 可以根据需要实现结构体为空的判断逻辑
		return false
	default:
		return false
	}
}

// getEmptyContainerValue 获取空容器的字符串表示
func getEmptyContainerValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		return "[]"
	case reflect.Map:
		return "{}"
	case reflect.Struct:
		return "{}"
	default:
		return ""
	}
}

// 更准确的显示宽度计算函数
func getDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if isWideChar(r) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

// isWideChar 判断是否为宽字符（简化版）
func isWideChar(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) ||
		(r >= 0x2E80 && r <= 0x9FFF) ||
		(r >= 0xAC00 && r <= 0xD7AF) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE30 && r <= 0xFE4F) ||
		(r >= 0xFF00 && r <= 0xFFEF)
}

// generateCompactStyleField 生成紧凑风格字段
func generateCompactStyleField(result *strings.Builder, field FieldInfo, indentStr string, options *Options) error {
	// 处理复杂类型（有子字段的情况）
	if field.HasChildren {
		if field.Comment != "" {
			// 看看子元素是否为空

			switch field.Field.Kind() {
			case reflect.Slice, reflect.Array:
				if field.Field.Len() == 0 {
					result.WriteString(fmt.Sprintf("%s%s: [] # %s\n", indentStr, field.Name, field.Comment))
				}
				return nil
			case reflect.Map:
				if field.Field.Len() == 0 {
					result.WriteString(fmt.Sprintf("%s%s: {} # %s\n", indentStr, field.Name, field.Comment))
				}
				return nil
			case reflect.Struct:
				fields := collectFieldInfo(field.Field, field.Field.Type(), field.FieldPath, options)
				if len(fields) == 0 {
					result.WriteString(fmt.Sprintf("%s%s: {} # %s\n", indentStr, field.Name, field.Comment))
				}
				return nil
			}

			result.WriteString(fmt.Sprintf("%s%s:  # %s", indentStr, field.Name, field.Comment))

		} else {
			result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))
		}
		err := generateFieldValue(result, field, indentStr, options)
		if err != nil {
			return err
		}
		return nil
	}

	// 处理数组/切片类型
	indent := 0
	hasVisibleChildren := false
	if field.Field.Kind() == reflect.Slice || field.Field.Kind() == reflect.Array {
		hasVisibleChildren = field.Field.Len() > 0
		if hasVisibleChildren {
			if field.Comment != "" {
				result.WriteString(fmt.Sprintf("%s%s: # %s", indentStr, field.Name, field.Comment))
			} else {
				result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))
			}
			indent = 1
		} else {
			result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))
		}
	} else {
		result.WriteString(fmt.Sprintf("%s%s: ", indentStr, field.Name))
	}

	// 生成字段值
	fieldValue, err := generateValue(field.Field, field.FieldPath, indent, options)
	if err != nil {
		return err
	}
	fieldValue = strings.TrimRight(fieldValue, "\n")

	// 输出最终结果
	if field.Comment != "" && !hasVisibleChildren {
		result.WriteString(fmt.Sprintf("%s # %s\n", fieldValue, field.Comment))
	} else {
		result.WriteString(fmt.Sprintf("%s\n", fieldValue))
	}

	return nil
}

// generateMinimalStyleField 生成最小风格字段
func generateMinimalStyleField(v interface{}) (string, error) {
	//yaml 直接转field.Field 成yaml
	yamlData, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// generateVerboseStyleField 生成详细风格字段
func generateVerboseStyleField(result *strings.Builder, field FieldInfo, indentStr string, options *Options) error {
	if field.Comment != "" {
		fieldTypeStr := field.Field.Type().String()
		result.WriteString(fmt.Sprintf("%s# %s (%s)\n", indentStr, field.Comment, fieldTypeStr))
	}
	result.WriteString(fmt.Sprintf("%s%s:", indentStr, field.Name))

	return generateFieldValue(result, field, indentStr, options)
}

// generateFieldValue 生成字段值
func generateFieldValue(result *strings.Builder, field FieldInfo, indentStr string, options *Options) error {
	// 特殊处理切片类型，即使它们没有复杂的子元素
	if field.HasChildren || field.Field.Kind() == reflect.Slice || field.Field.Kind() == reflect.Array {
		//如果元素和数组为空就不需要换行
		hasVisibleChildren := field.HasChildren ||
			(field.Field.Kind() == reflect.Slice && field.Field.Len() > 0) ||
			(field.Field.Kind() == reflect.Array && field.Field.Len() > 0)
		if hasVisibleChildren {
			result.WriteString("\n")
		}
		fieldValue, err := generateValue(field.Field, field.FieldPath, getIndentLevel(indentStr)+1, options)
		if err != nil {
			return err
		}
		if field.Field.Kind() == reflect.Slice || field.Field.Kind() == reflect.Array {
			lines := strings.Split(fieldValue, "\n")
			var commentLines, fieldLines []string
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					if strings.HasPrefix(strings.TrimSpace(line), "#") {
						commentLines = append(commentLines, fmt.Sprintf("%s\n", line))
					} else {
						fieldLines = append(fieldLines, normalizeTrailingNewlines1(fmt.Sprintf("%s\n", line)))
					}
				}
			}
			if len(commentLines) > 0 {
				result.WriteString(normalizeTrailingNewlines1(strings.Join(commentLines, "")))
			}
			if len(fieldLines) > 0 {
				result.WriteString(normalizeTrailingNewlines1(strings.Join(fieldLines, "")))
			}
		} else {
			result.WriteString(fieldValue)
		}
	} else {
		fieldValue, err := generateValue(field.Field, field.FieldPath, 0, options)
		if err != nil {
			return err
		}
		result.WriteString(" " + strings.TrimSpace(fieldValue))
	}
	return nil
}
func normalizeTrailingNewlines1(content string) string {
	// 如果内容以多个换行符结尾，移除所有尾部的换行符并转为1个
	if strings.HasSuffix(content, "\n") {
		trimmed := strings.TrimRight(content, "\n")
		return trimmed + "\n"
	}
	return content
}

// getIndentLevel 获取缩进级别
func getIndentLevel(indentStr string) int {
	return len(indentStr) / 2
}

// isValidKeyName 验证键名是否符合YAML标准
func isValidKeyName(key string) bool {
	if key == "" {
		return false
	}

	// 检查是否以数字、连字符或特殊字符开头
	if len(key) > 0 {
		firstChar := key[0]
		if (firstChar >= '0' && firstChar <= '9') || firstChar == '-' || firstChar == '_' {
			return false
		}
	}

	// 检查是否包含YAML特殊字符
	if strings.ContainsAny(key, ":#{}[]&*!|>\"'%?@`") {
		return false
	}

	// 检查是否为YAML关键字
	yamlKeywords := map[string]bool{
		"true": true, "false": true, "yes": true, "no": true,
		"on": true, "off": true, "null": true, "nil": true,
		"~": true, "True": true, "False": true, "TRUE": true, "FALSE": true,
		"Yes": true, "No": true, "YES": true, "NO": true,
		"On": true, "Off": true, "ON": true, "OFF": true,
		"Null": true, "NULL": true, "Nil": true, "NIL": true,
	}

	if yamlKeywords[key] {
		return false
	}

	return true
}

func getFieldName(fieldType reflect.StructField) string {
	// 检查yaml标签
	if yamlTag := fieldType.Tag.Get("yaml"); yamlTag != "" {
		if yamlTag == "-" {
			return "-"
		}
		parts := strings.Split(yamlTag, ",")
		if parts[0] != "" && parts[0] != "-" && !strings.Contains(parts[0], "=") && isValidKeyName(parts[0]) {
			return parts[0]
		}
	}

	// 检查yamlc标签
	if yamlcTag := fieldType.Tag.Get("yamlc"); yamlcTag != "" {
		if yamlcTag == "-" {
			return "-"
		}
		parts := strings.Split(yamlcTag, ",")
		if parts[0] != "" && parts[0] != "-" && !strings.Contains(parts[0], "=") && isValidKeyName(parts[0]) {
			return parts[0]
		}
	}

	// 如果没有标签，返回小写的字段名
	return "-"
}

// buildFieldPath 构建字段路径
func buildFieldPath(fieldPath, fieldName string) string {
	if fieldPath != "" {
		return fieldPath + "." + fieldName
	}
	return fieldName
}

// generateMap 生成Map YAML
func generateMap(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	if val.Len() == 0 {
		return " {}", nil
	}

	var result strings.Builder
	indentStr := strings.Repeat("  ", indent)

	iter := val.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		keyStr := fmt.Sprintf("%v", key.Interface())
		if needsQuoting(keyStr) {
			keyStr = fmt.Sprintf("%q", keyStr)
		}

		result.WriteString(fmt.Sprintf("%s%s:", indentStr, keyStr))

		if hasChildren(value) {
			result.WriteString("\n")
			valueStr, err := generateValue(value, fieldPath, indent+1, options)
			if err != nil {
				return "", err
			}
			result.WriteString(valueStr)
		} else {
			valueStr, err := generateValue(value, fieldPath, indent+1, options)
			if err != nil {
				return "", err
			}
			result.WriteString(" " + strings.TrimSpace(valueStr) + "\n")
		}
	}

	return result.String(), nil
}

// generateSlice 生成Slice YAML
func generateSlice(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	if val.Len() == 0 {
		return " []\n", nil
	}

	var result strings.Builder

	indentStr := strings.Repeat("  ", indent)

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)

		if hasChildren(item) {
			// 对于结构体等复杂类型，生成值并添加 "-" 前缀
			itemStr, err := generateValue(item, fieldPath, indent+1, options)
			if err != nil {
				return "", err
			}

			// 第一个元素保留注释，其他元素去掉注释
			keepComments := (i == 0)
			formattedStr := addDashPrefix(itemStr, indentStr, keepComments, options)
			result.WriteString(formattedStr)

			// 最后一个元素后添加换行
			if i == val.Len()-1 {
				result.WriteString("\n")
			}
		} else {
			if i == 0 {
				result.WriteString("\n")
			}
			// 简单类型，直接生成带 "- " 前缀的值
			itemStr, err := generateValue(item, fieldPath, indent+1, options)
			if err != nil {
				return "", err
			}

			trimmedValue := strings.TrimSpace(itemStr)
			if trimmedValue != "" {
				result.WriteString(fmt.Sprintf("%s- %s\n", indentStr, trimmedValue))
			} else {
				result.WriteString(fmt.Sprintf("%s-\n", indentStr))
			}
		}
	}

	return result.String(), nil
}

// addDashPrefix 为YAML列表项添加 "- " 前缀
func addDashPrefix(content string, indentStr string, keepComments bool, options *Options) string {
	lines := strings.Split(content, "\n")

	if !keepComments {
		// 去掉注释行，只保留非空非注释行
		var filteredLines []string
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
				filteredLines = append(filteredLines, line)
			}
		}
		lines = filteredLines
	}

	// 找到第一个非空非注释行，为其添加 "- " 前缀
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
			lines[i] = fmt.Sprintf("%s- %s", indentStr, trimmedLine)
			break
		}
	}

	return strings.Join(lines, "\n")
}

// generateString 生成字符串YAML
func generateString(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	str := val.String()

	// 验证字符串内容
	if err := validateStringContent(str); err != nil {
		return "", fmt.Errorf("invalid string content: %w", err)
	}

	if needsQuoting(str) {
		return fmt.Sprintf("%q", str), nil
	}
	return str, nil
}

// validateStringContent 验证字符串内容
func validateStringContent(str string) error {
	// 检查是否包含控制字符（除了常见的换行、制表符等）
	for _, r := range str {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			return fmt.Errorf("string contains invalid control character: %U", r)
		}
	}
	return nil
}

// generateInt 生成整数YAML
func generateInt(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	intVal := val.Int()

	// 验证整数范围
	switch val.Kind() {
	case reflect.Int8:
		if intVal < -128 || intVal > 127 {
			return "", fmt.Errorf("int8 value out of range: %d", intVal)
		}
	case reflect.Int16:
		if intVal < -32768 || intVal > 32767 {
			return "", fmt.Errorf("int16 value out of range: %d", intVal)
		}
	case reflect.Int32:
		if intVal < -2147483648 || intVal > 2147483647 {
			return "", fmt.Errorf("int32 value out of range: %d", intVal)
		}
	}

	return fmt.Sprintf("%d", intVal), nil
}

// generateUint 生成无符号整数YAML
func generateUint(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	uintVal := val.Uint()

	// 验证无符号整数范围
	switch val.Kind() {
	case reflect.Uint8:
		if uintVal > 255 {
			return "", fmt.Errorf("uint8 value out of range: %d", uintVal)
		}
	case reflect.Uint16:
		if uintVal > 65535 {
			return "", fmt.Errorf("uint16 value out of range: %d", uintVal)
		}
	case reflect.Uint32:
		if uintVal > 4294967295 {
			return "", fmt.Errorf("uint32 value out of range: %d", uintVal)
		}
	}

	return fmt.Sprintf("%d", uintVal), nil
}

// generateFloat 生成浮点数YAML
func generateFloat(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	floatVal := val.Float()

	// 验证浮点数有效性
	if isInvalidFloat(floatVal) {
		return "", fmt.Errorf("invalid float value: %f", floatVal)
	}

	// 根据类型确定精度
	switch val.Kind() {
	case reflect.Float32:
		return fmt.Sprintf("%.7g", float32(floatVal)), nil
	case reflect.Float64:
		return fmt.Sprintf("%.15g", floatVal), nil
	}

	return fmt.Sprintf("%g", floatVal), nil
}

// isInvalidFloat 检查浮点数是否有效
func isInvalidFloat(f float64) bool {
	// 检查是否为无穷大或NaN
	if f != f { // NaN检查
		return true
	}
	if f > 1.7976931348623157e+308 || f < -1.7976931348623157e+308 { // 无穷大检查
		return true
	}
	return false
}

// generateBool 生成布尔值YAML
func generateBool(val reflect.Value, fieldPath string, indent int, options *Options) (string, error) {
	return fmt.Sprintf("%t", val.Bool()), nil
}

// getComment 获取字段注释
func getComment(field reflect.StructField, fieldPath string, options *Options) string {
	// 1. 优先检查配置中的预设注释
	for _, commentMap := range options.Comments {
		if comment, exists := commentMap[fieldPath]; exists {
			return sanitizeComment(comment)
		}
	}

	// 2. 检查yamlc标签中的注释
	if yamlcTag := field.Tag.Get("yamlc"); yamlcTag != "" {
		parts := strings.Split(yamlcTag, ",")
		for _, part := range parts {
			if strings.HasPrefix(part, "comment=") {
				return sanitizeComment(strings.TrimPrefix(part, "comment="))
			}
		}
	}

	// 3. 检查comment标签
	if comment := field.Tag.Get("comment"); comment != "" {
		return sanitizeComment(comment)
	}

	// 4. 检查yaml标签中的注释
	if yamlTag := field.Tag.Get("yaml"); yamlTag != "" {
		parts := strings.Split(yamlTag, ",")
		for _, part := range parts {
			if strings.HasPrefix(part, "comment=") {
				return sanitizeComment(strings.TrimPrefix(part, "comment="))
			}
		}
	}

	return ""
}

// sanitizeComment 清理注释内容
func sanitizeComment(comment string) string {
	// 移除注释中的换行符和制表符，替换为空格
	comment = strings.ReplaceAll(comment, "\n", " ")
	comment = strings.ReplaceAll(comment, "\t", " ")
	comment = strings.ReplaceAll(comment, "\r", " ")

	// 移除多余的空格
	words := strings.Fields(comment)
	return strings.Join(words, " ")
}

// hasChildren 检查值是否有子元素
func hasChildren(val reflect.Value) bool {
	if !val.IsValid() {
		return false
	}

	switch val.Kind() {
	case reflect.Struct:
		return val.NumField() > 0
	case reflect.Map:
		return val.Len() > 0
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return false
		}
		// 检查第一个元素是否为复杂类型
		if val.Len() > 0 {
			firstItem := val.Index(0)
			return isComplexType(firstItem)
		}
		return false
	case reflect.Ptr:
		if val.IsNil() {
			return false
		}
		return hasChildren(val.Elem())
	case reflect.Interface:
		if val.IsNil() {
			return false
		}
		return hasChildren(val.Elem())
	default:
		return false
	}
}

// isComplexType 检查是否为复杂类型
func isComplexType(val reflect.Value) bool {
	if !val.IsValid() {
		return false
	}

	switch val.Kind() {
	case reflect.Struct, reflect.Map:
		return true
	case reflect.Slice, reflect.Array:
		return val.Len() > 0
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return false
		}
		return isComplexType(val.Elem())
	default:
		return false
	}
}

// needsQuoting 检查字符串是否需要引号
func needsQuoting(str string) bool {
	if str == "" {
		return true
	}

	// YAML保留字
	yamlKeywords := map[string]bool{
		"true": true, "false": true, "yes": true, "no": true,
		"on": true, "off": true, "null": true, "nil": true,
		"~": true, "True": true, "False": true, "TRUE": true, "FALSE": true,
		"Yes": true, "No": true, "YES": true, "NO": true,
		"On": true, "Off": true, "ON": true, "OFF": true,
		"Null": true, "NULL": true, "Nil": true, "NIL": true,
	}

	if yamlKeywords[str] {
		return true
	}

	// 检查是否为数字格式
	if isNumericString(str) {
		return true
	}

	// 检查特殊字符
	if strings.ContainsAny(str, ":#{}[]&*!|>\"'%?@`") {
		return true
	}

	// 检查换行符
	if strings.Contains(str, "\n") {
		return true
	}

	// 检查前后空格
	if strings.HasPrefix(str, " ") || strings.HasSuffix(str, " ") {
		return true
	}

	// 检查是否以数字开头的特殊格式
	if len(str) > 0 && (str[0] == '+' || str[0] == '-') {
		return true
	}

	return false
}

// isNumericString 检查字符串是否为数字格式
func isNumericString(str string) bool {
	if str == "" {
		return false
	}

	// 简单的数字格式检查
	for i, r := range str {
		if i == 0 && (r == '+' || r == '-') {
			continue
		}
		if !unicode.IsDigit(r) && r != '.' && r != 'e' && r != 'E' {
			return false
		}
	}

	// 更严格的检查：尝试解析为数字
	return strings.ContainsAny(str, "0123456789")
}

// SetGlobalStyle 设置全局注释风格
func SetGlobalStyle(style CommentStyle) {
	GlobalCommentStyle = style
}

// GetStyleFromString 从字符串获取风格枚举
func GetStyleFromString(styleStr string) CommentStyle {
	switch strings.ToLower(styleStr) {
	case "top":
		return StyleTop
	case "inline":
		return StyleInline
	case "smart":
		return StyleSmart
	case "compact":
		return StyleCompact
	case "minimal":
		return StyleMinimal
	case "verbose":
		return StyleVerbose
	case "spaced":
		return StyleSpaced
	case "grouped":
		return StyleGrouped
	case "sectioned":
		return StyleSectioned
	case "doc":
		return StyleDoc
	case "separate":
		return StyleSeparate
	default:
		return StyleSmart
	}
}

// ValidateOptions 验证选项配置
func ValidateOptions(options *Options) error {
	if options == nil {
		return fmt.Errorf("options cannot be nil")
	}

	// 验证注释风格范围
	if int(options.Style) < 0 || int(options.Style) > int(StyleSeparate) {
		return fmt.Errorf("invalid comment style: %d", options.Style)
	}

	// 验证注释内容
	for i, commentMap := range options.Comments {
		if commentMap == nil {
			return fmt.Errorf("comment map at index %d cannot be nil", i)
		}

		for fieldPath, comment := range commentMap {
			if fieldPath == "" {
				return fmt.Errorf("field path cannot be empty in comment map at index %d", i)
			}

			if err := validateCommentContent(comment); err != nil {
				return fmt.Errorf("invalid comment for field %q: %w", fieldPath, err)
			}
		}
	}

	return nil
}

// validateCommentContent 验证注释内容
func validateCommentContent(comment string) error {
	if len(comment) > 1000 {
		return fmt.Errorf("comment too long: %d characters (max 1000)", len(comment))
	}

	// 检查是否包含不合适的字符
	for _, r := range comment {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			return fmt.Errorf("comment contains invalid control character: %U", r)
		}
	}

	return nil
}

// GenWithValidation 生成YAML内容并进行全面验证
func GenWithValidation(v interface{}, opts ...Option) ([]byte, error) {
	// 预验证输入
	if v == nil {
		return nil, fmt.Errorf("input value cannot be nil")
	}

	// 验证反射值
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil, fmt.Errorf("input pointer cannot be nil")
	}

	// 构建和验证选项
	options := &Options{
		Style:    GlobalCommentStyle,
		Comments: make([]map[string]string, 0),
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := ValidateOptions(options); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	// 生成内容
	data, err := Gen(v, opts...)
	if err != nil {
		return nil, err
	}

	// 额外的后处理验证
	if err := ValidateGeneratedContent(data); err != nil {
		return nil, fmt.Errorf("generated content validation failed: %w", err)
	}

	return data, nil
}

// ValidateGeneratedContent 验证生成的内容
func ValidateGeneratedContent(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("generated content is empty")
	}

	// 检查是否有非UTF-8字符
	if !isValidUTF8(data) {
		return fmt.Errorf("generated content contains invalid UTF-8 sequences")
	}

	// 检查行结束符一致性
	content := string(data)
	if strings.Contains(content, "\r\n") && strings.Contains(content, "\n") {
		// 混合行结束符，标准化为\n
		content = strings.ReplaceAll(content, "\r\n", "\n")
		copy(data, []byte(content))
	}

	return nil
}

// isValidUTF8 检查数据是否为有效的UTF-8编码
func isValidUTF8(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			return false
		}
		data = data[size:]
	}
	return true
}
