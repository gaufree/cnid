// Package cnid 提供中国居民身份证和外国人永久居留身份证的校验、解析和生成功能
// 支持新旧版身份证，包括：
// - 中国居民身份证（15 位旧版和 18 位新版）
// - 外国人永久居留身份证（旧版 15 位和 2023 新版 18 位）
//
// MIT License
// Copyright (c) 2026 Mark Chen (gaufree)
package cnid

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 身份证类型常量
const (
	TypeUnknown         = iota // 未知类型
	TypeResidentOld15          // 中国居民身份证旧版（15 位）
	TypeResidentNew18          // 中国居民身份证新版（18 位）
	TypeForeignOld15           // 外国人永久居留身份证旧版（15 位，3 字母 +12 数字）
	TypeForeignNew18           // 外国人永久居留身份证新版（18 位，以 9 开头）
)

// 校验码系数（ISO 7064:1983.MOD 11-2）
var weightFactor = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// 校验码映射表
var checkCodeMap = []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

// 地区代码正则（简化的验证）
var areaCodePattern = regexp.MustCompile(`^[1-9]\d{5}$`)

// 外国人永久居留身份证新版正则（18 位，以 9 开头，第 2-4 位为字母国籍代码）
// 格式：9 + 国籍代码 (3 位字母) + 申领地代码 (2 位数字) + 出生日期 (8 位数字) + 顺序码 (3 位数字) + 校验码 (1 位)
var foreignNewPattern = regexp.MustCompile(`^9[A-Z]{3}\d{13}[0-9X]$`)

// 外国人永久居留身份证旧版正则（3 字母 +12 数字）
var foreignOldPattern = regexp.MustCompile(`^[A-Z]{3}\d{12}$`)

// IDInfo 身份证信息结构
type IDInfo struct {
	Type        int       // 身份证类型
	IDNumber    string    // 身份证号码（标准化后的大写格式）
	BirthDate   time.Time // 出生日期
	Gender      string    // 性别（"男" 或 "女"）
	AreaCode    string    // 地区代码
	Nationality string    // 国籍代码（仅外国人永久居留身份证）
	IssuePlace  string    // 申领地代码（仅外国人永久居留身份证）
}

// Validate 验证身份证号码是否有效
// 支持中国居民身份证（15 位/18 位）和外国人永久居留身份证（15 位/18 位）
func Validate(idNumber string) bool {
	if idNumber == "" {
		return false
	}

	idNumber = strings.ToUpper(strings.TrimSpace(idNumber))

	// 判断身份证类型
	idType := GetType(idNumber)
	if idType == TypeUnknown {
		return false
	}

	switch idType {
	case TypeResidentOld15:
		return validateResidentOld15(idNumber)
	case TypeResidentNew18:
		return validateResidentNew18(idNumber)
	case TypeForeignOld15:
		return validateForeignOld15(idNumber)
	case TypeForeignNew18:
		return validateForeignNew18(idNumber)
	default:
		return false
	}
}

// GetType 获取身份证类型
func GetType(idNumber string) int {
	if idNumber == "" {
		return TypeUnknown
	}

	idNumber = strings.ToUpper(strings.TrimSpace(idNumber))
	length := len(idNumber)

	// 18 位身份证
	if length == 18 {
		// 外国人永久居留身份证新版（以 9 开头）
		if strings.HasPrefix(idNumber, "9") && foreignNewPattern.MatchString(idNumber) {
			return TypeForeignNew18
		}
		// 中国居民身份证新版
		if residentNewPattern.MatchString(idNumber) {
			return TypeResidentNew18
		}
	}

	// 15 位身份证
	if length == 15 {
		// 外国人永久居留身份证旧版（3 字母 +12 数字）
		if foreignOldPattern.MatchString(idNumber) {
			return TypeForeignOld15
		}
		// 中国居民身份证旧版（纯数字）
		if residentOldPattern.MatchString(idNumber) {
			return TypeResidentOld15
		}
	}

	return TypeUnknown
}

// 中国居民身份证旧版正则（15 位纯数字）
var residentOldPattern = regexp.MustCompile(`^[1-9]\d{14}$`)

// 中国居民身份证新版正则（18 位，前 17 位数字，最后 1 位数字或 X）
var residentNewPattern = regexp.MustCompile(`^[1-9]\d{16}[0-9X]$`)

// validateResidentOld15 验证 15 位中国居民身份证
func validateResidentOld15(idNumber string) bool {
	if !residentOldPattern.MatchString(idNumber) {
		return false
	}

	// 验证地区代码
	areaCode := idNumber[:6]
	if !areaCodePattern.MatchString(areaCode) {
		return false
	}

	// 验证出生日期（YYMMDD 格式）
	yearStr := "19" + idNumber[6:8]
	monthStr := idNumber[8:10]
	dayStr := idNumber[10:12]

	return isValidDate(yearStr, monthStr, dayStr)
}

// validateResidentNew18 验证 18 位中国居民身份证
func validateResidentNew18(idNumber string) bool {
	if !residentNewPattern.MatchString(idNumber) {
		return false
	}

	// 验证地区代码
	areaCode := idNumber[:6]
	if !areaCodePattern.MatchString(areaCode) {
		return false
	}

	// 验证出生日期（YYYYMMDD 格式）
	yearStr := idNumber[6:10]
	monthStr := idNumber[10:12]
	dayStr := idNumber[12:14]

	if !isValidDate(yearStr, monthStr, dayStr) {
		return false
	}

	// 验证校验码
	return validateCheckCode(idNumber)
}

// validateForeignOld15 验证旧版外国人永久居留身份证（3 字母 +12 数字）
func validateForeignOld15(idNumber string) bool {
	// 先转大写再验证
	idUpper := strings.ToUpper(idNumber)
	return foreignOldPattern.MatchString(idUpper)
}

// validateForeignNew18 验证新版外国人永久居留身份证（18 位，以 9 开头）
func validateForeignNew18(idNumber string) bool {
	if !foreignNewPattern.MatchString(idNumber) {
		return false
	}

	// 验证校验码（外国人永居证需要特殊处理字母）
	return validateForeignCheckCode(idNumber)
}

// validateForeignCheckCode 验证外国人永久居留身份证的校验码
func validateForeignCheckCode(idNumber string) bool {
	body := idNumber[:17]
	expectedCheckCode := checkCodeMap[calculateForeignCheckCodeIndex(body)]
	actualCheckCode := string(idNumber[17])
	return actualCheckCode == expectedCheckCode
}

// calculateForeignCheckCodeIndex 计算外国人永久居留身份证校验码索引
func calculateForeignCheckCodeIndex(body string) int {
	sum := 0
	for i := 0; i < 17; i++ {
		var num int
		c := body[i]
		if c >= 'A' && c <= 'Z' {
			num = int(c - 'A' + 10)
		} else {
			num, _ = strconv.Atoi(string(c))
		}
		sum += num * weightFactor[i]
	}
	return sum % 11
}

// validateCheckCode 验证校验码（适用于 18 位身份证）
func validateCheckCode(idNumber string) bool {
	sum := 0
	for i := 0; i < 17; i++ {
		num, err := strconv.Atoi(string(idNumber[i]))
		if err != nil {
			return false
		}
		sum += num * weightFactor[i]
	}

	checkCodeIndex := sum % 11
	expectedCheckCode := checkCodeMap[checkCodeIndex]
	actualCheckCode := string(idNumber[17])

	return actualCheckCode == expectedCheckCode
}

// isValidDate 验证日期是否有效
func isValidDate(yearStr, monthStr, dayStr string) bool {
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > time.Now().Year() {
		return false
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return false
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > 31 {
		return false
	}

	// 验证具体日期的有效性（如 2 月 30 日等）
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date.Year() == year && int(date.Month()) == month && date.Day() == day
}

// Parse 解析身份证号码，返回详细信息
func Parse(idNumber string) (*IDInfo, error) {
	if idNumber == "" {
		return nil, fmt.Errorf("身份证号码不能为空")
	}

	idNumber = strings.ToUpper(strings.TrimSpace(idNumber))

	idType := GetType(idNumber)
	if idType == TypeUnknown {
		return nil, fmt.Errorf("无效的身份证号码格式")
	}

	info := &IDInfo{
		Type:     idType,
		IDNumber: idNumber,
	}

	var err error
	switch idType {
	case TypeResidentOld15:
		err = parseResidentOld15(idNumber, info)
	case TypeResidentNew18:
		err = parseResidentNew18(idNumber, info)
	case TypeForeignOld15:
		err = parseForeignOld15(idNumber, info)
	case TypeForeignNew18:
		err = parseForeignNew18(idNumber, info)
	}

	if err != nil {
		return nil, err
	}

	return info, nil
}

// parseResidentOld15 解析 15 位中国居民身份证
func parseResidentOld15(idNumber string, info *IDInfo) error {
	// 地区代码
	info.AreaCode = idNumber[:6]

	// 出生日期（YYMMDD，年份默认为 19xx）
	year := "19" + idNumber[6:8]
	month := idNumber[8:10]
	day := idNumber[10:12]

	birthDate, err := parseDate(year, month, day)
	if err != nil {
		return err
	}
	info.BirthDate = birthDate

	// 性别（第 15 位，奇数为男，偶数为女）
	genderCode, _ := strconv.Atoi(string(idNumber[14]))
	if genderCode%2 == 1 {
		info.Gender = "男"
	} else {
		info.Gender = "女"
	}

	return nil
}

// parseResidentNew18 解析 18 位中国居民身份证
func parseResidentNew18(idNumber string, info *IDInfo) error {
	// 地区代码
	info.AreaCode = idNumber[:6]

	// 出生日期（YYYYMMDD）
	year := idNumber[6:10]
	month := idNumber[10:12]
	day := idNumber[12:14]

	birthDate, err := parseDate(year, month, day)
	if err != nil {
		return err
	}
	info.BirthDate = birthDate

	// 性别（第 17 位，奇数为男，偶数为女）
	genderCode, _ := strconv.Atoi(string(idNumber[16]))
	if genderCode%2 == 1 {
		info.Gender = "男"
	} else {
		info.Gender = "女"
	}

	return nil
}

// parseForeignOld15 解析旧版外国人永久居留身份证
func parseForeignOld15(idNumber string, info *IDInfo) error {
	// 前 3 位为国籍代码
	info.Nationality = idNumber[:3]

	// 第 4-9 位为出生日期（YYYYMM）
	year := idNumber[3:7]
	month := idNumber[7:9]
	day := "01" // 旧版只有年月，日默认为 01

	birthDate, err := parseDate(year, month, day)
	if err != nil {
		return err
	}
	info.BirthDate = birthDate

	// 性别信息在第 10-12 位中编码，这里简化处理
	info.Gender = "未知"

	return nil
}

// parseForeignNew18 解析新版外国人永久居留身份证
func parseForeignNew18(idNumber string, info *IDInfo) error {
	// 第 1 位：外国人标识码（固定为 9）
	// 第 2-4 位：国籍代码
	info.Nationality = idNumber[1:4]

	// 第 5-6 位：申领地代码
	info.IssuePlace = idNumber[4:6]

	// 第 7-14 位：出生日期（YYYYMMDD）
	year := idNumber[6:10]
	month := idNumber[10:12]
	day := idNumber[12:14]

	birthDate, err := parseDate(year, month, day)
	if err != nil {
		return err
	}
	info.BirthDate = birthDate

	// 第 15-17 位：顺序码，其中第 17 位奇数为男，偶数为女
	genderCode, _ := strconv.Atoi(string(idNumber[16]))
	if genderCode%2 == 1 {
		info.Gender = "男"
	} else {
		info.Gender = "女"
	}

	return nil
}

// parseDate 解析日期字符串
func parseDate(yearStr, monthStr, dayStr string) (time.Time, error) {
	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)
	day, _ := strconv.Atoi(dayStr)

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

// GenerateResident 生成随机的中国居民身份证号码（18 位）
// areaCode: 6 位地区代码（可选，为空则随机生成）
// birthDate: 出生日期（可选，为空则随机生成 1950-2000 年之间的日期）
// gender: 性别（可选，"男"、"女" 或空表示随机）
func GenerateResident(areaCode string, birthDate time.Time, gender string) string {
	// 生成或验证地区代码
	if areaCode == "" {
		areaCode = generateAreaCode()
	} else if len(areaCode) != 6 || !areaCodePattern.MatchString(areaCode) {
		areaCode = generateAreaCode()
	}

	// 生成或验证出生日期
	var year, month, day int
	if birthDate.IsZero() {
		// 随机生成 1950-2000 年之间的日期
		year = randomInt(1950, 2000)
		month = randomInt(1, 12)
		day = randomInt(1, 28) // 简化处理，避免月末问题
	} else {
		year = birthDate.Year()
		month = int(birthDate.Month())
		day = birthDate.Day()
	}

	// 格式化出生日期
	dateStr := fmt.Sprintf("%04d%02d%02d", year, month, day)

	// 生成顺序码（第 15-17 位）
	var orderCode int
	if gender == "男" {
		orderCode = randomOdd(1, 999)
	} else if gender == "女" {
		orderCode = randomEven(1, 998)
	} else {
		orderCode = randomInt(1, 999)
	}
	orderStr := fmt.Sprintf("%03d", orderCode)

	// 组合前 17 位
	body := areaCode + dateStr + orderStr

	// 计算校验码
	checkCode := calculateCheckCode(body)

	return body + checkCode
}

// GenerateForeignNew 生成随机的新版外国人永久居留身份证号码（18 位）
// nationality: 3 位国籍代码（可选）
// issuePlace: 2 位申领地代码（可选）
// birthDate: 出生日期（可选）
// gender: 性别（可选）
func GenerateForeignNew(nationality, issuePlace string, birthDate time.Time, gender string) string {
	// 外国人标识码（固定为 9）
	prefix := "9"

	// 生成国籍代码
	if nationality == "" || len(nationality) != 3 {
		nationality = generateNationalityCode()
	}

	// 生成申领地代码
	if issuePlace == "" || len(issuePlace) != 2 {
		issuePlace = fmt.Sprintf("%02d", randomInt(11, 99))
	}

	// 生成出生日期
	var year, month, day int
	if birthDate.IsZero() {
		year = randomInt(1950, 2000)
		month = randomInt(1, 12)
		day = randomInt(1, 28)
	} else {
		year = birthDate.Year()
		month = int(birthDate.Month())
		day = birthDate.Day()
	}
	dateStr := fmt.Sprintf("%04d%02d%02d", year, month, day)

	// 生成顺序码
	var orderCode int
	if gender == "男" {
		orderCode = randomOdd(1, 999)
	} else if gender == "女" {
		orderCode = randomEven(1, 998)
	} else {
		orderCode = randomInt(1, 999)
	}
	orderStr := fmt.Sprintf("%03d", orderCode)

	// 组合前 17 位
	body := prefix + nationality + issuePlace + dateStr + orderStr

	// 计算校验码（外国人永居证需要特殊处理字母）
	checkCode := calculateForeignCheckCode(body)

	return body + checkCode
}

// calculateCheckCode 计算校验码（仅适用于纯数字的 17 位 body）
func calculateCheckCode(body string) string {
	sum := 0
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(string(body[i]))
		sum += num * weightFactor[i]
	}
	return checkCodeMap[sum%11]
}

// calculateForeignCheckCode 计算外国人永久居留身份证的校验码
// 格式：9 + 国籍代码 (3 位字母) + 申领地代码 (2 位数字) + 出生日期 (8 位数字) + 顺序码 (3 位数字)
// 需要将字母转换为数字：A=10, B=11, ..., Z=35
func calculateForeignCheckCode(body string) string {
	sum := 0
	for i := 0; i < 17; i++ {
		var num int
		c := body[i]
		if c >= 'A' && c <= 'Z' {
			num = int(c - 'A' + 10)
		} else {
			num, _ = strconv.Atoi(string(c))
		}
		sum += num * weightFactor[i]
	}
	return checkCodeMap[sum%11]
}

// generateAreaCode 生成随机的地区代码
func generateAreaCode() string {
	// 常见的省级行政区代码前缀
	provinceCodes := []int{11, 12, 13, 14, 15, 21, 22, 23, 31, 32, 33, 34, 35, 36, 37, 41, 42, 43, 44, 45, 46, 50, 51, 52, 53, 54, 61, 62, 63, 64, 65}
	province := provinceCodes[randomInt(0, len(provinceCodes)-1)]
	city := randomInt(1, 20)
	district := randomInt(1, 30)
	return fmt.Sprintf("%02d%02d%02d", province, city, district)
}

// generateNationalityCode 生成随机的国籍代码（3 位字母）
func generateNationalityCode() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 3)
	for i := range result {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[idx.Int64()]
	}
	return string(result)
}

// randomInt 生成指定范围内的随机整数
func randomInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

// randomOdd 生成指定范围内的随机奇数
func randomOdd(min, max int) int {
	num := randomInt(min, max)
	if num%2 == 0 {
		num++
	}
	if num > max {
		num -= 2
	}
	return num
}

// randomEven 生成指定范围内的随机偶数
func randomEven(min, max int) int {
	num := randomInt(min, max)
	if num%2 == 1 {
		num++
	}
	if num > max {
		num -= 2
	}
	return num
}

// UpgradeOld15To18 将 15 位中国居民身份证升级到 18 位
func UpgradeOld15To18(id15 string) (string, error) {
	if len(id15) != 15 || !residentOldPattern.MatchString(id15) {
		return "", fmt.Errorf("无效的 15 位身份证号码")
	}

	// 在年份前添加"19"
	body := id15[:6] + "19" + id15[6:]

	// 计算校验码
	checkCode := calculateCheckCode(body)

	return body + checkCode, nil
}

// GetTypeName 获取身份证类型的中文名称
func GetTypeName(idType int) string {
	switch idType {
	case TypeResidentOld15:
		return "中国居民身份证旧版（15 位）"
	case TypeResidentNew18:
		return "中国居民身份证新版（18 位）"
	case TypeForeignOld15:
		return "外国人永久居留身份证旧版（15 位）"
	case TypeForeignNew18:
		return "外国人永久居留身份证新版（18 位）"
	default:
		return "未知类型"
	}
}

// String IDInfo 的字符串表示
func (info *IDInfo) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("身份证类型：%s\n", GetTypeName(info.Type)))
	sb.WriteString(fmt.Sprintf("身份证号码：%s\n", info.IDNumber))
	sb.WriteString(fmt.Sprintf("出生日期：%s\n", info.BirthDate.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("性别：%s\n", info.Gender))

	if info.AreaCode != "" {
		sb.WriteString(fmt.Sprintf("地区代码：%s\n", info.AreaCode))
	}
	if info.Nationality != "" {
		sb.WriteString(fmt.Sprintf("国籍代码：%s\n", info.Nationality))
	}
	if info.IssuePlace != "" {
		sb.WriteString(fmt.Sprintf("申领地代码：%s\n", info.IssuePlace))
	}

	return sb.String()
}
