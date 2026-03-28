# CNID

中国居民身份证和外国人永久居留身份证 Golang 库。支持新旧版身份证的校验、解析和生成功能。

## 功能特性

- ✅ **中国居民身份证**
  - 支持 15 位旧版身份证
  - 支持 18 位新版身份证
  - 校验码验证（ISO 7064:1983.MOD 11-2）
  - 出生日期和性别解析
  - 15 位升级至 18 位

- ✅ **外国人永久居留身份证**
  - 支持旧版（15 位，3 字母 +12 数字）
  - 支持 2023 新版（18 位，以 9 开头，"五星卡"）
  - 校验码验证
  - 国籍代码、申领地代码解析
  - 出生日期和性别解析

- ✅ **其他功能**
  - 随机生成有效身份证号码
  - 指定参数生成（地区、出生日期、性别）
  - 大小写不敏感
  - 自动处理空白字符

## 安装

```bash
go get github.com/gaufree/cnid
```

## 快速开始

### 验证身份证号码

```go
package main

import (
    "fmt"
    "github.com/gaufree/cnid"
)

func main() {
    // 中国居民身份证新版（18 位）
    fmt.Println(cnid.Validate("11010519491231002X")) // true
    
    // 中国居民身份证旧版（15 位）
    fmt.Println(cnid.Validate("110105491231002")) // true
    
    // 外国人永久居留身份证新版（18 位）
    fmt.Println(cnid.Validate("9USA11199001018999")) // true
    
    // 外国人永久居留身份证旧版（15 位）
    fmt.Println(cnid.Validate("USA199001011234")) // true
}
```

### 获取身份证类型

```go
package main

import (
    "fmt"
    "github.com/gaufree/cnid"
)

func main() {
    idType := cnid.GetType("11010519491231002X")
    fmt.Println(cnid.GetTypeName(idType)) // 输出：中国居民身份证新版（18 位）
    
    idType = cnid.GetType("9USA11199001018999")
    fmt.Println(cnid.GetTypeName(idType)) // 输出：外国人永久居留身份证新版（18 位）
}
```

### 解析身份证信息

```go
package main

import (
    "fmt"
    "time"
    "github.com/gaufree/cnid"
)

func main() {
    // 解析中国居民身份证
    info, err := cnid.Parse("110105199001011234")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Printf("类型：%s\n", cnid.GetTypeName(info.Type))
    fmt.Printf("出生日期：%s\n", info.BirthDate.Format("2006-01-02"))
    fmt.Printf("性别：%s\n", info.Gender)
    fmt.Printf("地区代码：%s\n", info.AreaCode)
    
    // 解析外国人永久居留身份证
    foreignID := cnid.GenerateForeignNew("USA", "11", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "男")
    info, err = cnid.Parse(foreignID)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Printf("国籍代码：%s\n", info.Nationality)
    fmt.Printf("申领地代码：%s\n", info.IssuePlace)
}
```

### 生成身份证号码

```go
package main

import (
    "fmt"
    "time"
    "github.com/gaufree/cnid"
)

func main() {
    // 随机生成中国居民身份证
    id := cnid.GenerateResident("", time.Time{}, "")
    fmt.Println(id)
    
    // 指定参数生成
    id = cnid.GenerateResident("110105", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "男")
    fmt.Println(id)
    
    // 随机生成外国人永久居留身份证（新版）
    foreignID := cnid.GenerateForeignNew("", "", time.Time{}, "")
    fmt.Println(foreignID)
    
    // 指定参数生成
    foreignID = cnid.GenerateForeignNew("USA", "11", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "女")
    fmt.Println(foreignID)
}
```

### 15 位身份证升级至 18 位

```go
package main

import (
    "fmt"
    "github.com/gaufree/cnid"
)

func main() {
    id18, err := cnid.UpgradeOld15To18("110105491231002")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(id18) // 输出：11010519491231002X
}
```

## API 文档

### 常量

```go
const (
    TypeUnknown        = iota // 未知类型
    TypeResidentOld15         // 中国居民身份证旧版（15 位）
    TypeResidentNew18         // 中国居民身份证新版（18 位）
    TypeForeignOld15          // 外国人永久居留身份证旧版（15 位）
    TypeForeignNew18          // 外国人永久居留身份证新版（18 位）
)
```

### 函数

#### `Validate(idNumber string) bool`
验证身份证号码是否有效。

#### `GetType(idNumber string) int`
获取身份证类型。

#### `GetTypeName(idType int) string`
获取身份证类型的中文名称。

#### `Parse(idNumber string) (*IDInfo, error)`
解析身份证号码，返回详细信息。

#### `GenerateResident(areaCode string, birthDate time.Time, gender string) string`
生成随机的中国居民身份证号码（18 位）。
- `areaCode`: 6 位地区代码（可选，为空则随机生成）
- `birthDate`: 出生日期（可选，为空则随机生成 1950-2000 年之间的日期）
- `gender`: 性别（"男"、"女" 或空表示随机）

#### `GenerateForeignNew(nationality, issuePlace string, birthDate time.Time, gender string) string`
生成随机的新版外国人永久居留身份证号码（18 位）。
- `nationality`: 3 位国籍代码（可选）
- `issuePlace`: 2 位申领地代码（可选）
- `birthDate`: 出生日期（可选）
- `gender`: 性别（"男"、"女" 或空表示随机）

#### `UpgradeOld15To18(id15 string) (string, error)`
将 15 位中国居民身份证升级到 18 位。

### 类型

#### `IDInfo`
身份证信息结构。

```go
type IDInfo struct {
    Type        int       // 身份证类型
    IDNumber    string    // 身份证号码（标准化后的大写格式）
    BirthDate   time.Time // 出生日期
    Gender      string    // 性别（"男" 或 "女"）
    AreaCode    string    // 地区代码
    Nationality string    // 国籍代码（仅外国人永久居留身份证）
    IssuePlace  string    // 申领地代码（仅外国人永久居留身份证）
}
```

## 身份证编码规则

### 中国居民身份证（18 位）
```
110105 19491231 002 X
  |        |      |   |
  |        |      |   └─ 校验码
  |        |      └───── 顺序码（奇数男，偶数女）
  |        └──────────── 出生日期（YYYYMMDD）
  └───────────────────── 地区代码
```

### 中国居民身份证（15 位旧版）
```
110105 491231 002
  |       |     |
  |       |     └─ 顺序码
  |       └─────── 出生日期（YYMMDD）
  └─────────────── 地区代码
```

### 外国人永久居留身份证（18 位新版/"五星卡"）
```
9 USA 11 19900101 123 X
|  |    |     |      |   |
|  |    |     |      |   └─ 校验码
|  |    |     |      └───── 顺序码（奇数男，偶数女）
|  |    |     └──────────── 出生日期（YYYYMMDD）
|  |    └────────────────── 申领地代码
|  └─────────────────────── 国籍代码（3 位字母）
└────────────────────────── 外国人标识码（固定为 9）
```

### 外国人永久居留身份证（15 位旧版）
```
USA 19900101 1234
 |     |       |
 |     |       └─ 顺序码
 |     └───────── 出生日期（YYYYMMDD，实际只有年月）
 └─────────────── 国籍代码（3 位字母）
```

## 运行测试

```bash
go test -v
```

## 许可证

MIT License

Copyright (c) 2026 Mark Chen (gaufree)

## 免责声明

本库仅供学习和测试使用，不得用于非法用途。使用本库生成的身份证号码仅用于测试目的，请勿用于真实场景。
