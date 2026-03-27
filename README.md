# CNID - 中国身份证 Golang 库

支持新版外国人居留身份证的中国身份证验证库。

## 功能特性

- ✅ 支持中国公民身份证 (18位) 解析与验证
- ✅ 支持外国人居留身份证 (16位) 解析与验证
- ✅ 自动识别身份证类型
- ✅ 校验码验证
- ✅ 地区代码解析
- ✅ 国籍解析 (外国人居留身份证)
- ✅ 出生日期解析
- ✅ 性别解析
- ✅ 年龄计算
- ✅ 脱敏处理 (Mask/Hide)
- ✅ 完整的单元测试

## 安装

```bash
go get github.com/gaufree/cnid
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/gaufree/cnid"
)

func main() {
    // 解析中国公民身份证
    id, err := cnid.ParseChinese("110101199001010017")
    if err != nil {
        fmt.Println("错误:", err)
        return
    }
    fmt.Printf("地区: %s\n", cnid.GetRegion(id.Number))
    fmt.Printf("出生日期: %s\n", id.Birthday.Format("2006-01-02"))
    fmt.Printf("性别: %d (1=男,0=女)\n", id.Sex)

    // 解析外国人居留身份证 (16位)
    fid, err := cnid.ParseForeigner("8011100900101013")
    if err != nil {
        fmt.Println("错误:", err)
        return
    }
    fmt.Printf("国籍: %s\n", fid.CountryName)
    fmt.Printf("地区: %s\n", cnid.GetRegion(fid.Number))

    // 自动识别类型
    info, err := cnid.Parse("110101199001010017")
    if err != nil {
        fmt.Println("错误:", err)
        return
    }
    switch info.GetType() {
    case cnid.TypeChinese:
        fmt.Println("中国公民身份证")
    case cnid.TypeForeigner:
        fmt.Println("外国人居留身份证")
    }
}
```

## 身份证格式

### 中国公民身份证 (18位)

```
[地区代码:6位] + [出生日期:8位] + [顺序码:3位] + [校验码:1位]
110101        + 19900101       + 123            + X
```

### 外国人居留身份证 (16位)

```
8 + [国籍代码:2位] + [地区代码:4位] + [出生日期:6位] + [顺序码:2位] + [校验码:1位]
8 + 01            + 1100          + 900101         + 01             + 3
```

**格式说明:**
- 第1位: 固定为 `8`
- 第2-3位: 国籍代码 (01-99, 数字代码)
- 第4-7位: 停留/居留地地区代码 (4位)
- 第8-13位: 出生日期 (YYMMDD, 两位年)
- 第14-15位: 顺序码 (奇数=男,偶数=女)
- 第16位: 校验码

## API 文档

### 核心函数

| 函数 | 描述 |
|------|------|
| `Parse(idNumber string) (IDInfo, error)` | 自动识别并解析身份证 |
| `ParseChinese(idNumber string) (*ChineseID, error)` | 解析中国公民身份证 |
| `ParseForeigner(idNumber string) (*ForeignerID, error)` | 解析外国人居留身份证 |
| `Validate(idNumber string) bool` | 验证身份证是否有效 |
| `GetIDType(idNumber string) IDType` | 获取身份证类型 |

### 辅助函数

| 函数 | 描述 |
|------|------|
| `GetRegion(idNumber string) string` | 获取地区名称 |
| `GetCountry(idNumber string) string` | 获取国籍 |
| `GetAge(idNumber string) (int, error)` | 计算年龄 |
| `GetBirthdayString(idNumber string) (string, error)` | 获取格式化出生日期 |
| `GetSexString(idNumber string) (string, error)` | 获取性别字符串 |
| `Mask(idNumber string) string` | 脱敏处理 (显示前后部分) |
| `Hide(idNumber string) string` | 隐藏出生日期 |

### 数据结构

#### ChineseID

```go
type ChineseID struct {
    Number     string    // 身份证号码
    RegionCode string    // 地区代码
    Birthday   time.Time // 出生日期
    Sex        int       // 性别: 0-女, 1-男
    Sequence   int       // 顺序码
    Checksum   rune      // 校验码
}
```

#### ForeignerID

```go
type ForeignerID struct {
    Number      string    // 证件号码
    CountryCode string    // 国籍代码 (2位数字)
    CountryName string    // 国籍名称
    RegionCode  string    // 地区代码
    Birthday    time.Time // 出生日期
    Sex         int       // 性别: 0-女, 1-男
    Sequence    int       // 顺序码
    IsPermanent bool      // 是否长期有效
    Checksum    rune      // 校验码
}
```

### IDType 常量

```go
const (
    TypeUnknown   IDType = 0  // 未知类型
    TypeChinese   IDType = 1  // 中国公民身份证
    TypeForeigner IDType = 2  // 外国人居留身份证
)
```

## 运行测试

```bash
go test -v ./...
```

## 运行基准测试

```bash
go test -bench=. -benchmem
```

## 许可证

MIT License - © gaufree