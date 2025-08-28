# YAMLC - Go语言YAML生成器（带注释支持）

一个功能强大的Go语言YAML生成器，支持多种注释风格和灵活的配置选项。

## 特性

- 🎨 **多种注释风格**: 支持10种不同的注释风格
- 🏷️ **灵活的注释方式**: 支持结构体标签和自定义注释映射
- 🔄 **智能格式化**: 根据数据结构自动选择最佳注释位置
- 📁 **文件操作**: 直接写入文件或io.Writer
- ✅ **YAML验证**: 内置YAML格式验证
- 🚀 **高性能**: 使用反射和缓冲区优化性能

## 安装

```bash
go get github.com/binrclab/yamlc
```

## 支持的注释风格

根据测试结果，YAMLC支持以下注释风格：

### 1. **StyleTop** (默认)
注释显示在字段上方：
```yaml
# 用户姓名
name: 张三
# 用户年龄
age: 30
# 电子邮箱地址
email: "zhangsan@example.com"
```

### 2. **StyleInline**
注释显示在行内，完美对齐：
```yaml
name: 张三                                     # 用户姓名
age: 30                                        # 用户年龄
email: "zhangsan@example.com"                  # 电子邮箱地址
phone: 13800138000                             # 电话号码
address:                                       # 用户地址
  street: 中关村大街1号                   # 街道地址
  city: 北京                              # 城市
  country: 中国                          # 国家
  zipcode: 100080                         # 邮政编码
```

### 3. **StyleSmart**
智能选择：简单字段使用行内注释，复杂结构使用顶部注释：
```yaml
name: 张三                                     # 用户姓名
age: 30                                        # 用户年龄
email: "zhangsan@example.com"                  # 电子邮箱地址
phone: 13800138000                             # 电话号码
# 用户地址
address: 
  street: 中关村大街1号                   # 街道地址
  city: 北京                              # 城市
  country: 中国                          # 国家
  zipcode: 100080                         # 邮政编码
```

### 4. **StyleCompact**
紧凑风格，最小间距：
```yaml
name: 张三 # 用户姓名
age: 30 # 用户年龄
email: "zhangsan@example.com" # 电子邮箱地址
phone: 13800138000 # 电话号码
address: # 用户地址
  street: 中关村大街1号 # 街道地址
  city: 北京 # 城市
  country: 中国 # 国家
  zipcode: 100080 # 邮政编码
```

### 5. **StyleMinimal**
无注释，纯净的YAML输出：
```yaml
name: 张三
age: 30
email: zhangsan@example.com
phone: 13800138000
address:
    street: 中关村大街1号
    city: 北京
    country: 中国
    zipcode: 100080
tags:
    - 开发者
    - Go语言
    - 后端
```

### 6. **StyleVerbose**
详细风格，包含类型信息：
```yaml
# 用户姓名 (string)
name: 张三
# 用户年龄 (int)
age: 30
# 电子邮箱地址 (string)
email: "zhangsan@example.com"
# 电话号码 (int64)
phone: 13800138000
# 用户地址 (*yamlc.Address)
address:
  # 街道地址 (string)
  street: 中关村大街1号
  # 城市 (string)
  city: 北京
```

### 7. **StyleSpaced**
字段间增加空行：
```yaml
# 用户姓名
name: 张三

# 用户年龄
age: 30

# 电子邮箱地址
email: "zhangsan@example.com"

# 电话号码
phone: 13800138000

# 用户地址
address: 
  # 街道地址
  street: 中关村大街1号
  # 城市
  city: 北京
```

### 8. **StyleGrouped**
逻辑分组风格：
```yaml
# 用户姓名
name: 张三
# 用户年龄
age: 30
# 电子邮箱地址
email: "zhangsan@example.com"
# 电话号码
phone: 13800138000

# 用户地址
address: 
  # 街道地址
  street: 中关村大街1号
  # 城市
  city: 北京
  # 国家
  country: 中国
  # 邮政编码
  zipcode: 100080
```

### 9. **StyleSectioned**
分节风格，注释集中在每节开头：
```yaml
# 用户姓名
# 用户年龄
# 电子邮箱地址
# 电话号码

name: 张三
age: 30
email: "zhangsan@example.com"
phone: 13800138000

# 用户地址
address: 
  # 街道地址
  # 城市
  # 国家
  # 邮政编码

  street: 中关村大街1号
  city: 北京
  country: 中国
  zipcode: 100080
```

### 10. **StyleDoc** 和 **StyleSeparate**
文档风格，带头部注释块：
```yaml
############################################
# name(string):用户姓名
# age(int):用户年龄
# email(string):电子邮箱地址
# phone(int64):电话号码
# address(*yamlc.Address):用户地址
###########################################

name: 张三
age: 30
email: "zhangsan@example.com"
phone: 13800138000
address: 
  ############################################
  # street(string):街道地址
  # city(string):城市
  # country(string):国家
  # zipcode(int):邮政编码
  ###########################################

  street: 中关村大街1号
  city: 北京
  country: 中国
  zipcode: 100080
```

## 基本使用

### 定义结构体

```go
type User struct {
    Name     string   `yaml:"name" comment:"用户姓名"`
    Age      int      `yaml:"age" comment:"用户年龄"`
    Email    string   `yaml:"email" comment:"电子邮箱地址"`
    Phone    int64    `yaml:"phone" comment:"电话号码"`
    Active   bool     `yaml:"active" comment:"是否激活"`
    Tags     []string `yaml:"tags" comment:"用户标签"`
    Address  *Address `yaml:"address" comment:"用户地址"`
}

type Address struct {
    Street  string `yaml:"street" comment:"街道地址"`
    City    string `yaml:"city" comment:"城市"`
    Country string `yaml:"country" comment:"国家"`
    Zipcode int    `yaml:"zipcode" comment:"邮政编码"`
}
```

### 生成YAML

```go
package main

import (
    "fmt"
    "github.com/binrclab/yamlc"
)

func main() {
    user := User{
        Name:   "张三",
        Age:    30,
        Email:  "zhangsan@example.com",
        Phone:  13800138000,
        Active: true,
        Tags:   []string{"开发者", "Go语言", "后端"},
        Address: &Address{
            Street:  "中关村大街1号",
            City:    "北京", 
            Country: "中国",
            Zipcode: 100080,
        },
    }

    // 使用默认风格 (StyleTop)
    yaml, err := yamlc.Gen(user)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(yaml))
}
```

## 高级使用

### 使用不同的注释风格

```go
// 行内注释风格
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleInline))

// 智能风格
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleSmart))

// 紧凑风格
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleCompact))

// 最小风格（无注释）
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleMinimal))

// 详细风格（带类型信息）
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleVerbose))
```

### 自定义注释映射

```go
comments := map[string]string{
    "name":           "用户的完整姓名",
    "age":            "用户年龄（岁）",
    "email":          "主要联系邮箱",
    "address.street": "详细街道地址",
    "address.city":   "所在城市",
}

yaml, err := yamlc.Gen(user, 
    yamlc.WithComments(comments),
    yamlc.WithStyle(yamlc.StyleInline))
```

### 写入文件

```go
// 写入文件
err := yamlc.GenToFile(user, "config.yaml", 
    yamlc.WithStyle(yamlc.StyleInline))

// 写入io.Writer
var buf bytes.Buffer
err := yamlc.GenToWriter(user, &buf,
    yamlc.WithStyle(yamlc.StyleSmart))
```

### 验证选项

```go
// 生成时进行验证
yaml, err := yamlc.Gen(user,
    yamlc.WithStyle(yamlc.StyleInline),
    yamlc.WithValidation(true))

// 不进行验证（更快）
yaml, err := yamlc.Gen(user,
    yamlc.WithStyle(yamlc.StyleInline),
    yamlc.WithValidation(false))
```

## 配置选项

YAMLC提供多种配置选项：

```go
type Options struct {
    Style      CommentStyle         // 注释风格
    Comments   map[string]string    // 自定义注释映射
    Validation bool                 // 启用YAML验证
    Indent     string              // 缩进字符串（默认："  "）
    MaxWidth   int                 // 对齐的最大行宽
}
```

### 可用选项

- `WithStyle(style CommentStyle)` - 设置注释风格
- `WithComments(comments map[string]string)` - 设置自定义注释
- `WithValidation(enabled bool)` - 启用/禁用YAML验证
- `WithIndent(indent string)` - 设置自定义缩进
- `WithMaxWidth(width int)` - 设置对齐的最大行宽

## 测试结果示例

基于测试结果，以下是实际应用示例：

### 用户数据结构
```go
type User struct {
    Name           string             `comment:"用户姓名"`
    Age            int                `comment:"用户年龄"`
    Email          string             `comment:"电子邮箱地址"`
    Phone          int64              `comment:"电话号码"`
    Address        *Address           `comment:"用户地址"`
    Address2       *Address           `comment:"用户地址2"`
    Tags           []string           `comment:"用户标签"`
    Cog            []string           `comment:"用户标签"`
    Tesc           *Tesc2            `comment:"用户标签"`
    WorkExperience []*WorkExperience `comment:"任职经历"`
    Active         bool               `comment:"是否激活"`
    Score          float64            `comment:"用户评分"`
}

type WorkExperience struct {
    Company  string `comment:"公司名"`
    Position string `comment:"职位"`
}
```

### 不同风格输出效果

**StyleInline - 完美对齐：**
```yaml
name: 张三                                     # 用户姓名
age: 30                                        # 用户年龄
email: "zhangsan@example.com"                  # 电子邮箱地址
phone: 13800138000                             # 电话号码
workExperience:                                # 任职经历
  - company: co1                             # 公司名
    position: golang开发                     # 职位
  - company: co2                             # 公司名
    position: 高级golang开发                 # 职位
active: true                                   # 是否激活
score: 95.5                                    # 用户评分
```

**StyleSmart - 智能选择：**
```yaml
name: 张三                                     # 用户姓名
age: 30                                        # 用户年龄
email: "zhangsan@example.com"                  # 电子邮箱地址
phone: 13800138000                             # 电话号码
# 任职经历
workExperience: 
  - company: co1                             # 公司名
    position: golang开发                     # 职位
  - company: co2                             # 公司名
    position: 高级golang开发                 # 职位
active: true                                   # 是否激活
score: 95.5                                    # 用户评分
```

## 性能考虑

- 使用缓冲字符串构建优化性能
- 反射缓存避免重复类型分析
- 可配置验证平衡速度与准确性
- Unicode感知的文本宽度计算确保正确对齐

## 错误处理

```go
yaml, err := yamlc.Gen(data)
if err != nil {
    switch err {
    case yamlc.ErrInvalidData:
        log.Println("输入数据无效")
    case yamlc.ErrValidationFailed:
        log.Println("生成的YAML验证失败")
    default:
        log.Printf("意外错误: %v", err)
    }
}
```

## 特色功能

### Unicode支持
- 完美支持中文、日文、韩文等宽字符
- 正确计算字符显示宽度
- 注释对齐考虑字符宽度差异

### 空容器处理
- 智能处理空数组 `[]` 和空对象 `{}`
- 根据上下文选择合适的显示方式

### 灵活配置
- 支持结构体标签和映射注释
- 多种输出目标（字符串、文件、Writer）
- 可选的YAML格式验证

## 贡献

欢迎贡献代码！请随时提交Pull Request。

## 许可证

本项目采用MIT许可证 - 查看LICENSE文件了解详情。

## 更新日志

### v1.0.0
- 初始版本发布
- 支持10种不同的注释风格
- Unicode感知的文本对齐
- 灵活的配置选项
- 内置YAML验证
- 全面的测试覆盖