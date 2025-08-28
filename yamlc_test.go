package yamlc

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// 测试数据结构
type Address struct {
	Street  string `yaml:"street,omitempty"  yamlc:"comment=街道地址"`
	City    string `yaml:"city,omitempty"    yamlc:"comment=城市"`
	Country string `yaml:"country,omitempty" yamlc:"comment=国家"`
	Zipcode int    `yaml:"zipcode,omitempty" yamlc:"comment=邮政编码"`
}

type WorkExperience struct {
	Company  string `yaml:"company,omitempty"  yamlc:"comment=公司名"`
	Position string `yaml:"position,omitempty" yamlc:"comment=职位"`
}

type Tesc2 struct {
}

type User struct {
	Name           string            `yaml:"name,omitempty"           yamlc:"comment=用户姓名"`
	Age            int               `yaml:"age,omitempty"            yamlc:"comment=用户年龄"`
	Email          string            `yaml:"email,omitempty"          yamlc:"comment=电子邮箱地址"`
	Phone          int64             `yaml:"phone,omitempty"          yamlc:"comment=电话号码"`
	Address        *Address          `yaml:"address,omitempty"        yamlc:"comment=用户地址"`
	Address2       *Address          `yaml:"address2,omitempty"       yamlc:"comment=用户地址2"`
	Tags           []string          `yaml:"tags,omitempty"           yamlc:"comment=用户标签"`
	Cog            []string          `yaml:"cog,omitempty"            yamlc:"comment=用户标签"`
	Tesc           *Tesc2            `yaml:"tesc,omitempty"           yamlc:"comment=用户标签"`
	WorkExperience []*WorkExperience `yaml:"workExperience,omitempty" yamlc:"comment=任职经历"`
	Active         bool              `yaml:"active,omitempty"         yamlc:"comment=是否激活"`
	Score          float64           `yaml:"score,omitempty"          yamlc:"comment=用户评分"`
	PrivateField   string            // 未导出字段，应该被忽略
	IgnoredField   string            `yaml:"-"` // 被忽略的字段
}

// 创建测试用户数据
func createTestUser() *User {
	return &User{
		Name:   "张三",
		Age:    30,
		Email:  "zhangsan@example.com",
		Phone:  13800138000,
		Active: true,
		Score:  95.5,
		Address: &Address{
			Street:  "中关村大街1号",
			City:    "北京",
			Country: "中s国",
			Zipcode: 100080,
		},
		Address2: &Address{},
		Tags:     []string{"开发者", "Go语言", "后端"},
		Cog:      []string{},
		Tesc:     &Tesc2{},
		WorkExperience: []*WorkExperience{
			{
				Company:  "co1",
				Position: "golang开发",
			},
			{
				Company:  "co2",
				Position: "高级golang开发",
			},
		},
		PrivateField: "私有字段",
		IgnoredField: "被忽略的字段",
	}
}

func TestDebugGen(t *testing.T) {
	user := createTestUser()

	styles := []CommentStyle{
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

	for _, style := range styles {

		content, err := Gen(user, WithStyle(style), WithComment(map[string]string{
			// "name": "name",
		}))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("-----------------------%s-------------------\n", GetStyleString(int(style)))

		fmt.Printf("Generated YAML (without validation):\n%s\n", content)
		//检查是否是yaml文件格式
		err = Validate([]byte(content))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Generated YAML (without validation):\n%s\n", content)

		fmt.Printf("--------------------------------------------------------------------------\n")
	}
}

// 测试 Gen 函数
func TestGen(t *testing.T) {
	user := createTestUser()

	// 测试默认风格
	data, err := Gen(user)
	if err != nil {
		t.Fatalf("Gen failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Generated data is empty")
	}

	// 验证生成的YAML包含关键字段
	yamlStr := string(data)
	if !strings.Contains(yamlStr, "name: 张三") {
		t.Error("Generated YAML missing name field")
	}
	if !strings.Contains(yamlStr, "age: 30") {
		t.Error("Generated YAML missing age field")
	}
	if !strings.Contains(yamlStr, "address:") {
		t.Error("Generated YAML missing address field")
	}

	// 验证私有字段被忽略
	if strings.Contains(yamlStr, "PrivateField") {
		t.Error("Private field should not appear in generated YAML")
	}
	if strings.Contains(yamlStr, "IgnoredField") {
		t.Error("Ignored field should not appear in generated YAML")
	}
}

// 测试不同注释风格
func TestGenWithDifferentStyles(t *testing.T) {
	user := createTestUser()

	styles := []CommentStyle{
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

	for _, style := range styles {
		t.Run(GetStyleString(int(style)), func(t *testing.T) {
			data, err := Gen(user, WithStyle(style))
			if err != nil {
				t.Fatalf("Gen with style %s failed: %v", GetStyleString(int(style)), err)
			}

			if len(data) == 0 {
				t.Fatal("Generated data is empty")
			}

			// 验证风格特定的特征
			yamlStr := string(data)
			switch style {
			case StyleTop:
				if !strings.Contains(yamlStr, "# 用户姓名") {
					t.Error("Top style should have comment above field")
				}
			case StyleInline:
				if !strings.Contains(yamlStr, "# 用户姓名") {
					t.Error("Inline style should have inline comments")
				}
			case StyleMinimal:
				if strings.Contains(yamlStr, "#") {
					t.Error("Minimal style should not have comments")
				}
			case StyleVerbose:
				if !strings.Contains(yamlStr, "(string)") {
					t.Error("Verbose style should show type information")
				}
			}
		})
	}
}

// 测试 Write 函数
func TestWrite(t *testing.T) {
	user := createTestUser()
	var buf bytes.Buffer

	err := Write(&buf, user)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("Buffer is empty after write")
	}

	// 验证写入的内容
	content := buf.String()
	if !strings.Contains(content, "name: 张三") {
		t.Error("Written content missing name field")
	}
}

// 测试 WriteFile 函数
func TestWriteFile(t *testing.T) {
	user := createTestUser()
	filename := "test_user.yaml"

	// 清理测试文件
	defer os.Remove(filename)

	err := WriteFile(filename, user)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// 验证文件被创建
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("File was not created")
	}

	// 验证文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if !strings.Contains(string(content), "name: 张三") {
		t.Error("File content missing name field")
	}
}

// 测试验证函数
func TestValidateYAML(t *testing.T) {
	validYAML := []byte(`
name: 张三
age: 30
address:
  street: 中关村大街1号
  city: 北京
`)

	invalidYAML := []byte(`
name: 张三
  age: 30  # 错误的缩进
`)

	// 测试有效YAML
	err := ValidateYAML(validYAML)
	if err != nil {
		t.Errorf("ValidateYAML should pass for valid YAML: %v", err)
	}

	// 测试无效YAML
	err = ValidateYAML(invalidYAML)
	if err == nil {
		t.Error("ValidateYAML should fail for invalid YAML")
	}
}

// 测试结构验证
func TestValidateStructure(t *testing.T) {
	validYAML := []byte(`
name: 张三
age: 30
address:
  street: 中关村大街1号
  city: 北京
`)

	invalidYAML := []byte(`
name: 张三
 age: 30  # 奇数缩进
`)

	// 测试有效结构
	err := ValidateStructure(validYAML)
	if err != nil {
		t.Errorf("ValidateStructure should pass for valid structure: %v", err)
	}

	// 测试无效结构
	err = ValidateStructure(invalidYAML)
	if err == nil {
		t.Error("ValidateStructure should fail for invalid structure")
	}
}

// 测试选项函数
func TestOptions(t *testing.T) {
	user := createTestUser()

	// 测试自定义注释
	comments := map[string]string{
		"name":    "自定义姓名注释",
		"age":     "自定义年龄注释",
		"address": "自定义地址注释",
	}

	data, err := Gen(user, WithComment(comments))
	if err != nil {
		t.Fatalf("Gen with custom comments failed: %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "自定义姓名注释") {
		t.Error("Custom comment not found in generated YAML")
	}
}

// 测试边界情况
func TestEdgeCases(t *testing.T) {
	// 测试空值
	_, err := Gen(nil)
	if err == nil {
		t.Error("Gen should fail for nil input")
	}

	// 测试空指针
	var user *User
	_, err = Gen(user)
	if err == nil {
		t.Error("Gen should fail for nil pointer")
	}

	// 测试空结构体
	emptyUser := &User{}
	data, err := Gen(emptyUser)
	if err != nil {
		t.Fatalf("Gen should work for empty struct: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Generated data should not be empty for empty struct")
	}
}

// 测试复杂数据类型
func TestComplexTypes(t *testing.T) {
	complexUser := &User{
		Name: "测试用户",
		Tags: []string{"标签1", "标签2", "标签3"},
		WorkExperience: []*WorkExperience{
			{Company: "公司A", Position: "职位A"},
			{Company: "公司B", Position: "职位B"},
		},
	}

	data, err := Gen(complexUser)
	if err != nil {
		t.Fatalf("Gen failed for complex types: %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "- 标签1") {
		t.Error("Slice elements not properly generated")
	}
	if !strings.Contains(yamlStr, "- company: 公司A") {
		t.Error("Struct slice elements not properly generated")
	}
}

// 测试全局风格设置
func TestGlobalStyle(t *testing.T) {
	originalStyle := GetStyle()
	defer SetGlobalStyle(originalStyle)

	user := createTestUser()

	// 设置全局风格
	SetGlobalStyle(StyleVerbose)
	if GetStyle() != StyleVerbose {
		t.Error("Global style not set correctly")
	}

	// 测试全局风格生效
	data, err := Gen(user)
	if err != nil {
		t.Fatalf("Gen with global style failed: %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "(string)") {
		t.Error("Global verbose style not applied")
	}
}

// 测试风格字符串转换
func TestStyleStringConversion(t *testing.T) {
	testCases := []struct {
		style    int
		expected string
	}{
		{0, "top"},
		{1, "inline"},
		{2, "smart"},
		{3, "compact"},
		{4, "minimal"},
		{5, "verbose"},
		{6, "spaced"},
		{7, "grouped"},
		{8, "sectioned"},
		{9, "doc"},
		{10, "separate"},
		{999, "smart"}, // 默认值
	}

	for _, tc := range testCases {
		result := GetStyleString(tc.style)
		if result != tc.expected {
			t.Errorf("GetStyleString(%d) = %s, expected %s", tc.style, result, tc.expected)
		}
	}
}

// 测试从字符串获取风格
func TestGetStyleFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected CommentStyle
	}{
		{"top", StyleTop},
		{"inline", StyleInline},
		{"smart", StyleSmart},
		{"compact", StyleCompact},
		{"minimal", StyleMinimal},
		{"verbose", StyleVerbose},
		{"spaced", StyleSpaced},
		{"grouped", StyleGrouped},
		{"sectioned", StyleSectioned},
		{"doc", StyleDoc},
		{"separate", StyleSeparate},
		{"unknown", StyleSmart}, // 默认值
		{"TOP", StyleTop},       // 大小写不敏感
		{"Smart", StyleSmart},
	}

	for _, tc := range testCases {
		result := GetStyleFromString(tc.input)
		if result != tc.expected {
			t.Errorf("GetStyleFromString(%s) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

// 测试选项验证
func TestValidateOptions(t *testing.T) {
	// 测试有效选项
	validOptions := &Options{
		Style: StyleTop,
		Comments: []map[string]string{
			{"name": "测试注释"},
		},
	}

	err := ValidateOptions(validOptions)
	if err != nil {
		t.Errorf("ValidateOptions should pass for valid options: %v", err)
	}

	// 测试无效选项
	invalidOptions := &Options{
		Style: CommentStyle(999), // 无效风格
	}

	err = ValidateOptions(invalidOptions)
	if err == nil {
		t.Error("ValidateOptions should fail for invalid style")
	}

	// 测试空选项
	err = ValidateOptions(nil)
	if err == nil {
		t.Error("ValidateOptions should fail for nil options")
	}
}

// 测试带验证的生成
func TestGenWithValidation(t *testing.T) {
	user := createTestUser()

	data, err := GenWithValidation(user)
	if err != nil {
		t.Fatalf("GenWithValidation failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Generated data is empty")
	}

	// 验证生成的内容是有效的YAML
	err = ValidateYAML(data)
	if err != nil {
		t.Errorf("Generated content is not valid YAML: %v", err)
	}
}

// 测试字段信息收集
func TestCollectFieldInfo(t *testing.T) {
	user := createTestUser()
	val := reflect.ValueOf(user).Elem()
	typ := val.Type()

	fields := collectFieldInfo(val, typ, "", &Options{})
	if len(fields) == 0 {
		t.Fatal("No fields collected")
	}

	// 验证字段数量（应该排除私有字段和被忽略的字段）
	expectedFieldCount := 9 // 导出的字段数量
	if len(fields) != expectedFieldCount {
		t.Errorf("Expected %d fields, got %d", expectedFieldCount, len(fields))
	}

	// 验证字段信息
	foundName := false
	for _, field := range fields {
		if field.Name == "Name" {
			foundName = true
			if field.Comment != "用户姓名" {
				t.Errorf("Expected comment '用户姓名', got '%s'", field.Comment)
			}
			break
		}
	}

	if !foundName {
		t.Error("Name field not found in collected fields")
	}
}

// 测试字符串引号处理
func TestNeedsQuoting(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", true},             // 空字符串
		{"hello", false},       // 普通字符串
		{"true", true},         // YAML关键字
		{"123", true},          // 数字格式
		{"hello world", false}, // 包含空格
		{"hello:world", true},  // 包含特殊字符
		{"hello\nworld", true}, // 包含换行符
		{" hello", true},       // 前导空格
		{"hello ", true},       // 尾随空格
		{"+123", true},         // 以+开头
		{"-123", true},         // 以-开头
	}

	for _, tc := range testCases {
		result := needsQuoting(tc.input)
		if result != tc.expected {
			t.Errorf("needsQuoting(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

// 测试性能基准
func BenchmarkGen(b *testing.B) {
	user := createTestUser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Gen(user)
		if err != nil {
			b.Fatalf("Gen failed: %v", err)
		}
	}
}

func BenchmarkGenWithStyle(b *testing.B) {
	user := createTestUser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Gen(user, WithStyle(StyleVerbose))
		if err != nil {
			b.Fatalf("Gen with style failed: %v", err)
		}
	}
}

// 测试辅助函数
func TestHelperFunctions(t *testing.T) {
	// 测试 getFieldName
	fieldType := reflect.StructField{
		Name: "TestField",
		Tag:  `yaml:"custom_name"`,
	}

	fieldName := getFieldName(fieldType)
	if fieldName != "custom_name" {
		t.Errorf("Expected field name 'custom_name', got '%s'", fieldName)
	}

	// 测试 buildFieldPath
	path := buildFieldPath("parent", "child")
	if path != "parent.child" {
		t.Errorf("Expected path 'parent.child', got '%s'", path)
	}

	// 测试 getIndentLevel
	level := getIndentLevel("    ")
	if level != 2 {
		t.Errorf("Expected indent level 2, got %d", level)
	}
}
