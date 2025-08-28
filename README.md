# YAMLC - Go YAML Generator with Comment Support

A powerful Go YAML generator that supports multiple comment styles and flexible configuration options.

## Features

- 🎨 **Multiple Comment Styles**: Support for 10 different comment styles
- 🏷️ **Flexible Comment Options**: Support for struct tags and custom comment mappings
- 🔄 **Smart Formatting**: Automatically choose the best comment placement based on data structure
- 📁 **File Operations**: Direct writing to files or io.Writer
- ✅ **YAML Validation**: Built-in YAML format validation
- 🚀 **High Performance**: Optimized with reflection and buffer usage

## Installation

```bash
go get github.com/binrclab/yamlc
```

## Supported Comment Styles

Based on the test results, YAMLC supports the following comment styles:

### 1. **StyleTop** (Default)
Comments appear above fields:
```yaml
# User name
name: John Doe
# User age  
age: 30
# Email address
email: "john@example.com"
```

### 2. **StyleInline**
Comments appear inline with proper alignment:
```yaml
name: John Doe                                  # User name
age: 30                                         # User age
email: "john@example.com"                       # Email address
phone: 13800138000                              # Phone number
```

### 3. **StyleSmart**
Smart placement: inline for simple fields, top for complex structures:
```yaml
name: John Doe                                  # User name
age: 30                                         # User age
# User address
address: 
  street: Main Street 1                         # Street address
  city: New York                                # City
```

### 4. **StyleCompact**
Compact style with minimal spacing:
```yaml
name: John Doe # User name
age: 30 # User age
email: "john@example.com" # Email address
address: # User address
  street: Main Street 1 # Street address
  city: New York # City
```

### 5. **StyleMinimal**
No comments, clean YAML output:
```yaml
name: John Doe
age: 30
email: john@example.com
phone: 13800138000
address:
    street: Main Street 1
    city: New York
    country: USA
    zipcode: 10001
```

### 6. **StyleVerbose**
Detailed comments with type information:
```yaml
# User name (string)
name: John Doe
# User age (int)
age: 30
# Email address (string)
email: "john@example.com"
# Phone number (int64)
phone: 13800138000
```

### 7. **StyleSpaced**
Extra spacing between field groups:
```yaml
# User name
name: John Doe

# User age
age: 30

# Email address
email: "john@example.com"

# Phone number
phone: 13800138000
```

### 8. **StyleGrouped**
Logical grouping of related fields:
```yaml
# User name
name: John Doe
# User age
age: 30
# Email address
email: "john@example.com"
# Phone number
phone: 13800138000

# User address
address: 
  # Street address
  street: Main Street 1
  # City
  city: New York
```

### 9. **StyleSectioned**
Comments grouped at the beginning of each section:
```yaml
# User name
# User age
# Email address
# Phone number

name: John Doe
age: 30
email: "john@example.com"
phone: 13800138000

# User address
address: 
  # Street address
  # City
  # Country
  # Postal code

  street: Main Street 1
  city: New York
  country: USA
  zipcode: 10001
```

### 10. **StyleDoc** & **StyleSeparate**
Documentation-style with header blocks:
```yaml
############################################
# name(string): User name
# age(int): User age
# email(string): Email address
# phone(int64): Phone number
# address(*Address): User address
###########################################

name: John Doe
age: 30
email: "john@example.com"
phone: 13800138000
address: 
  street: Main Street 1
  city: New York
  country: USA
  zipcode: 10001
```

## Basic Usage

### Define Your Struct

```go
type User struct {
    Name     string   `yaml:"name" comment:"User name"`
    Age      int      `yaml:"age" comment:"User age"`
    Email    string   `yaml:"email" comment:"Email address"`
    Phone    int64    `yaml:"phone" comment:"Phone number"`
    Active   bool     `yaml:"active" comment:"Is user active"`
    Tags     []string `yaml:"tags" comment:"User tags"`
    Address  *Address `yaml:"address" comment:"User address"`
}

type Address struct {
    Street  string `yaml:"street" comment:"Street address"`
    City    string `yaml:"city" comment:"City"`
    Country string `yaml:"country" comment:"Country"`
    Zipcode int    `yaml:"zipcode" comment:"Postal code"`
}
```

### Generate YAML

```go
package main

import (
    "fmt"
    "github.com/binrclab/yamlc"
)

func main() {
    user := User{
        Name:   "John Doe",
        Age:    30,
        Email:  "john@example.com",
        Phone:  13800138000,
        Active: true,
        Tags:   []string{"developer", "golang", "backend"},
        Address: &Address{
            Street:  "Main Street 1",
            City:    "New York", 
            Country: "USA",
            Zipcode: 10001,
        },
    }

    // Use default style (StyleTop)
    yaml, err := yamlc.Gen(user)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(yaml))
}
```

## Advanced Usage

### Using Different Comment Styles

```go
// Inline comment style
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleInline))

// Smart style
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleSmart))

// Compact style
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleCompact))

// Minimal style (no comments)
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleMinimal))

// Verbose style with type information
yaml, err := yamlc.Gen(user, yamlc.WithStyle(yamlc.StyleVerbose))
```

### Custom Comment Mappings

```go
comments := map[string]string{
    "name":           "Full name of the user",
    "age":            "Age in years",
    "email":          "Primary email contact",
    "address.street": "Street address line",
    "address.city":   "City name",
}

yaml, err := yamlc.Gen(user, 
    yamlc.WithComments(comments),
    yamlc.WithStyle(yamlc.StyleInline))
```

### Write to File

```go
// Write to file
err := yamlc.GenToFile(user, "config.yaml", 
    yamlc.WithStyle(yamlc.StyleInline))

// Write to io.Writer
var buf bytes.Buffer
err := yamlc.GenToWriter(user, &buf,
    yamlc.WithStyle(yamlc.StyleSmart))
```

### Validation Options

```go
// Generate with validation
yaml, err := yamlc.Gen(user,
    yamlc.WithStyle(yamlc.StyleInline),
    yamlc.WithValidation(true))

// Generate without validation (faster)
yaml, err := yamlc.Gen(user,
    yamlc.WithStyle(yamlc.StyleInline),
    yamlc.WithValidation(false))
```

## Configuration Options

YAMLC provides various configuration options:

```go
type Options struct {
    Style      CommentStyle         // Comment style
    Comments   map[string]string    // Custom comment mappings
    Validation bool                 // Enable YAML validation
    Indent     string              // Indentation string (default: "  ")
    MaxWidth   int                 // Maximum line width for alignment
}
```

### Available Options

- `WithStyle(style CommentStyle)` - Set comment style
- `WithComments(comments map[string]string)` - Set custom comments
- `WithValidation(enabled bool)` - Enable/disable YAML validation
- `WithIndent(indent string)` - Set custom indentation
- `WithMaxWidth(width int)` - Set maximum line width for alignment

## Examples from Test Results

Based on the test results, here are practical examples:

### User Data Structure
```go
type User struct {
    Name           string             `comment:"User name"`
    Age            int                `comment:"User age"`
    Email          string             `comment:"Email address"`
    Phone          int64              `comment:"Phone number"`
    Address        *Address           `comment:"User address"`
    Address2       *Address           `comment:"User address 2"`
    Tags           []string           `comment:"User tags"`
    Cog            []string           `comment:"User tags"`
    Tesc           *Tesc2            `comment:"User tags"`
    WorkExperience []*WorkExperience `comment:"Work experience"`
    Active         bool               `comment:"Is active"`
    Score          float64            `comment:"User score"`
}
```

### Different Style Outputs

**StyleInline - Perfect Alignment:**
```yaml
name: 张三                                     # User name
age: 30                                        # User age
email: "zhangsan@example.com"                  # Email address
phone: 13800138000                             # Phone number
address:                                       # User address
  street: 中关村大街1号                   # Street address
  city: 北京                              # City
  country: 中国                          # Country
  zipcode: 100080                         # Postal code
```

**StyleSmart - Intelligent Placement:**
```yaml
name: 张三                                     # User name
age: 30                                        # User age
email: "zhangsan@example.com"                  # Email address
phone: 13800138000                             # Phone number
# User address
address: 
  street: 中关村大街1号                   # Street address
  city: 北京                              # City
  country: 中国                          # Country
  zipcode: 100080                         # Postal code
```

## Performance Considerations

- Uses buffered string building for optimal performance
- Reflection caching to avoid repeated type analysis
- Configurable validation to balance speed vs accuracy
- Unicode-aware text width calculation for proper alignment

## Error Handling

```go
yaml, err := yamlc.Gen(data)
if err != nil {
    switch err {
    case yamlc.ErrInvalidData:
        log.Println("Invalid input data")
    case yamlc.ErrValidationFailed:
        log.Println("Generated YAML failed validation")
    default:
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Changelog

### v1.0.0
- Initial release
- Support for 10 different comment styles
- Unicode-aware text alignment
- Flexible configuration options
- Built-in YAML validation
- Comprehensive test coverage