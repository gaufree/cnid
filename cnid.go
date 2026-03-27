// Package cnid provides Chinese ID card validation and parsing.
// MIT License - © gaufree
package cnid

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// IDType 身份证类型
type IDType int

const (
	TypeUnknown IDType = iota
	TypeChinese    // 中国公民身份证 (18位)
	TypeForeigner  // 外国人居留身份证 (18位)
)

// ErrInvalidID 无效的身份证号码
var ErrInvalidID = errors.New("invalid ID number")

// ChineseID 中国公民身份证结构
type ChineseID struct {
	Number     string    // 身份证号码
	RegionCode string    // 地区代码
	Birthday   time.Time // 出生日期
	Sex        int       // 性别: 0-女, 1-男
	Sequence   int       // 顺序码
	Checksum   rune      // 校验码
}

// ForeignerID 外国人居留身份证结构
type ForeignerID struct {
	Number      string    // 证件号码
	CountryCode string    // 国籍代码 (2位数字)
	CountryName string    // 国籍名称
	RegionCode  string    // 地区代码 (停留或居留地)
	Birthday    time.Time // 出生日期
	Sex         int       // 性别: 0-女, 1-男
	Sequence    int       // 顺序码
	IsPermanent bool      // 是否长期有效
	Checksum    rune      // 校验码
}

// IDInfo 统一的信息接口
type IDInfo interface {
	GetType() IDType
	GetNumber() string
	GetBirthday() time.Time
	GetSex() int
	GetRegionCode() string
	IsValid() bool
}

// GetType 获取身份证类型
func (id *ChineseID) GetType() IDType { return TypeChinese }

// GetNumber 获取身份证号码
func (id *ChineseID) GetNumber() string { return id.Number }

// GetBirthday 获取出生日期
func (id *ChineseID) GetBirthday() time.Time { return id.Birthday }

// GetSex 获取性别 (0:女, 1:男)
func (id *ChineseID) GetSex() int { return id.Sex }

// GetRegionCode 获取地区代码
func (id *ChineseID) GetRegionCode() string { return id.RegionCode }

// IsValid 验证身份证是否有效
func (id *ChineseID) IsValid() bool {
	return id.Number != "" && !id.Birthday.IsZero()
}

// GetType 获取身份证类型
func (id *ForeignerID) GetType() IDType { return TypeForeigner }

// GetNumber 获取证件号码
func (id *ForeignerID) GetNumber() string { return id.Number }

// GetBirthday 获取出生日期
func (id *ForeignerID) GetBirthday() time.Time { return id.Birthday }

// GetSex 获取性别 (0:女, 1:男)
func (id *ForeignerID) GetSex() int { return id.Sex }

// GetRegionCode 获取地区代码
func (id *ForeignerID) GetRegionCode() string { return id.RegionCode }

// IsValid 验证证件是否有效
func (id *ForeignerID) IsValid() bool {
	return id.Number != "" && !id.Birthday.IsZero()
}

// chineseIDRegex 中国公民身份证正则 (18位)
var chineseIDRegex = regexp.MustCompile(`^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)

// foreignerIDRegex 外国人居留身份证正则 (16位)
// 格式: 8 + 国家代码(2位) + 地区代码(4位) + 出生日期(6位YYMMDD) + 顺序码(2位) + 校验码(1位)
var foreignerIDRegex = regexp.MustCompile(`^8\d{2}\d{4}\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{2}[\dXx]$`)

// 权重系数 (用于校验码计算)
var weightFactors = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// 校验码映射
var checksumMap = map[int]rune{
	0: '0', 1: '1', 2: '2', 3: '3', 4: '4',
	5: '5', 6: '6', 7: '7', 8: '8', 9: '9',
	10: 'X',
}

// 2位数字国家代码映射 (部分常见国家)
var countryCodeMapping = map[string]string{
	"01": "美国",
	"02": "英国",
	"03": "德国",
	"04": "法国",
	"05": "日本",
	"06": "韩国",
	"07": "俄罗斯",
	"08": "加拿大",
	"09": "澳大利亚",
	"10": "印度",
	"11": "巴西",
	"12": "意大利",
	"13": "西班牙",
	"14": "墨西哥",
	"15": "荷兰",
	"16": "瑞典",
	"17": "挪威",
	"18": "丹麦",
	"19": "芬兰",
	"20": "瑞士",
	"21": "比利时",
	"22": "奥地利",
	"23": "波兰",
	"24": "乌克兰",
	"25": "哈萨克斯坦",
	"26": "新加坡",
	"27": "马来西亚",
	"28": "泰国",
	"29": "越南",
	"30": "菲律宾",
	"31": "印度尼西亚",
	"32": "巴基斯坦",
	"33": "阿联酋",
	"34": "沙特阿拉伯",
	"35": "埃及",
	"36": "尼日利亚",
	"37": "南非",
	"38": "阿根廷",
	"39": "智利",
	"40": "哥伦比亚",
	"41": "秘鲁",
	"42": "新西兰",
	"43": "葡萄牙",
	"44": "希腊",
	"45": "捷克",
	"46": "匈牙利",
	"47": "罗马尼亚",
	"48": "以色列",
	"49": "土耳其",
}

// 地区代码映射表(部分)
var regionCodeMapping = map[string]string{
	"110000": "北京市",
	"110101": "北京市东城区",
	"110102": "北京市西城区",
	"110105": "北京市朝阳区",
	"110106": "北京市丰台区",
	"110107": "北京市石景山区",
	"110108": "北京市海淀区",
	"310000": "上海市",
	"310101": "上海市黄浦区",
	"310104": "上海市徐汇区",
	"310105": "上海市长宁区",
	"310106": "上海市静安区",
	"310107": "上海市普陀区",
	"440000": "广东省",
	"440100": "广东省广州市",
	"440300": "广东省深圳市",
	"330000": "浙江省",
	"330100": "浙江省杭州市",
	"320000": "江苏省",
	"320100": "江苏省南京市",
	"500000": "重庆市",
	"510000": "四川省",
	"510100": "四川省成都市",
	"610000": "陕西省",
	"610100": "陕西省西安市",
}

// Parse 解析身份证号码,自动识别类型
func Parse(idNumber string) (IDInfo, error) {
	// 去除空格
	idNumber = cleanSpaces(idNumber)

	// 先尝试解析为中国公民身份证
	if chineseIDRegex.MatchString(idNumber) {
		return ParseChinese(idNumber)
	}

	// 再尝试解析为外国人居留身份证
	if foreignerIDRegex.MatchString(idNumber) {
		return ParseForeigner(idNumber)
	}

	return nil, ErrInvalidID
}

// ParseChinese 解析中国公民身份证
func ParseChinese(idNumber string) (*ChineseID, error) {
	idNumber = cleanSpaces(idNumber)

	if !chineseIDRegex.MatchString(idNumber) {
		return nil, ErrInvalidID
	}

	id := &ChineseID{Number: idNumber}

	// 解析地区代码
	id.RegionCode = idNumber[0:6]
	if _, ok := regionCodeMapping[id.RegionCode]; !ok {
		// 地区代码不在预定义表中,只截取前6位
	}

	// 解析出生日期
	birthYear, _ := strconv.Atoi(idNumber[6:10])
	birthMonth, _ := strconv.Atoi(idNumber[10:12])
	birthDay, _ := strconv.Atoi(idNumber[12:14])
	id.Birthday = time.Date(birthYear, time.Month(birthMonth), birthDay, 0, 0, 0, 0, time.UTC)

	// 解析顺序码和性别
	seq, _ := strconv.Atoi(idNumber[14:17])
	id.Sequence = seq
	id.Sex = seq % 2 // 奇数为男,偶数为女

	// 解析校验码
	id.Checksum = rune(idNumber[17])

	// 验证校验码
	if !validateChecksum(idNumber) {
		return nil, ErrInvalidID
	}

	return id, nil
}

// ParseForeigner 解析外国人居留身份证
func ParseForeigner(idNumber string) (*ForeignerID, error) {
	idNumber = cleanSpaces(idNumber)

	if !foreignerIDRegex.MatchString(idNumber) {
		return nil, ErrInvalidID
	}

	id := &ForeignerID{Number: idNumber}

	// 解析国家代码 (第2-3位)
	countryCode := idNumber[1:3]
	id.CountryCode = countryCode
	if name, ok := countryCodeMapping[countryCode]; ok {
		id.CountryName = name
	} else {
		id.CountryName = countryCode
	}

	// 解析地区代码 (第4-7位,4位)
	id.RegionCode = idNumber[3:7]
	if _, ok := regionCodeMapping[id.RegionCode]; !ok {
		// 地区代码不在预定义表中
	}

	// 解析出生日期 (第8-13位, 格式为YYMMDD)
	birthYY, _ := strconv.Atoi(idNumber[7:9])
	birthMM, _ := strconv.Atoi(idNumber[9:11])
	birthDD, _ := strconv.Atoi(idNumber[11:13])
	
	// 两位数年份需要转换为四位数 (假设1900-1999: 00-99 对应 1900-1999, 00-99 对应 2000-2099)
	birthYear := 1900 + birthYY
	if birthYY < 50 {
		birthYear = 2000 + birthYY
	}
	id.Birthday = time.Date(birthYear, time.Month(birthMM), birthDD, 0, 0, 0, 0, time.UTC)

	// 解析顺序码和性别 (第14-15位)
	seq, _ := strconv.Atoi(idNumber[13:15])
	id.Sequence = seq
	id.Sex = seq % 2 // 奇数为男,偶数为女

	// 解析校验码 (第16位)
	id.Checksum = rune(idNumber[15])

	// 验证校验码
	if !validateChecksum(idNumber) {
		return nil, ErrInvalidID
	}

	return id, nil
}

// validateChecksum 验证校验码
func validateChecksum(idNumber string) bool {
	if len(idNumber) != 18 && len(idNumber) != 16 {
		return false
	}

	// 根据长度选择要验证的位数
	checkLen := 15 // 对于16位证件,验证前15位
	if len(idNumber) == 18 {
		checkLen = 17
	}
	
	idToCheck := idNumber[:checkLen]
	
	sum := 0
	for i := 0; i < checkLen; i++ {
		digit, _ := strconv.Atoi(string(idToCheck[i]))
		sum += digit * weightFactors[i]
	}

	remainder := sum % 11
	expectedChecksum := checksumMap[remainder]

	return string(idNumber[checkLen]) == string(expectedChecksum) || idNumber[checkLen] == 'x' || idNumber[checkLen] == 'X'
}

// Validate 验证身份证号码是否有效(通用函数)
func Validate(idNumber string) bool {
	_, err := Parse(idNumber)
	return err == nil
}

// ValidateChinese 验证中国公民身份证号码是否有效
func ValidateChinese(idNumber string) bool {
	_, err := ParseChinese(idNumber)
	return err == nil
}

// ValidateForeigner 验证外国人居留身份证号码是否有效
func ValidateForeigner(idNumber string) bool {
	_, err := ParseForeigner(idNumber)
	return err == nil
}

// GetRegion 获取地区名称
func GetRegion(idNumber string) string {
	cleaned := cleanSpaces(idNumber)
	if len(cleaned) >= 6 {
		regionCode := cleaned[0:6]
		if name, ok := regionCodeMapping[regionCode]; ok {
			return name
		}
	}
	return ""
}

// GetCountry 获取国籍(仅对外国人居留身份证有效)
func GetCountry(idNumber string) string {
	cleaned := cleanSpaces(idNumber)
	if len(cleaned) >= 3 && cleaned[0] == '8' {
		countryCode := cleaned[1:3]
		if name, ok := countryCodeMapping[countryCode]; ok {
			return name
		}
		return "" // 未知国家代码返回空
	}
	return ""
}

// GetIDType 获取身份证类型
func GetIDType(idNumber string) IDType {
	cleaned := cleanSpaces(idNumber)
	if chineseIDRegex.MatchString(cleaned) {
		return TypeChinese
	}
	if foreignerIDRegex.MatchString(cleaned) {
		return TypeForeigner
	}
	return TypeUnknown
}

// GetAge 计算年龄
func GetAge(idNumber string) (int, error) {
	info, err := Parse(idNumber)
	if err != nil {
		return 0, err
	}

	birthday := info.GetBirthday()
	now := time.Now()
	age := now.Year() - birthday.Year()

	// 检查是否已经过了生日
	if now.YearDay() < birthday.YearDay() {
		age--
	}

	return age, nil
}

// GetBirthdayString 格式化输出出生日期
func GetBirthdayString(idNumber string) (string, error) {
	info, err := Parse(idNumber)
	if err != nil {
		return "", err
	}

	return info.GetBirthday().Format("2006-01-02"), nil
}

// GetSexString 获取性别字符串
func GetSexString(idNumber string) (string, error) {
	info, err := Parse(idNumber)
	if err != nil {
		return "", err
	}

	sex := info.GetSex()
	if sex == 1 {
		return "男", nil
	}
	return "女", nil
}

// String 实现Stringer接口
func (id *ChineseID) String() string {
	return fmt.Sprintf("ChineseID{Number: %s, Region: %s, Birthday: %s, Sex: %d}",
		id.Number, GetRegion(id.Number), id.Birthday.Format("2006-01-02"), id.Sex)
}

// String 实现Stringer接口
func (id *ForeignerID) String() string {
	return fmt.Sprintf("ForeignerID{Number: %s, Country: %s, Region: %s, Birthday: %s, Sex: %d}",
		id.Number, id.CountryName, GetRegion(id.Number), id.Birthday.Format("2006-01-02"), id.Sex)
}

// cleanSpaces 去除字符串中的空格
func cleanSpaces(s string) string {
	result := ""
	for _, c := range s {
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			result += string(c)
		}
	}
	return result
}

// Mask 脱敏处理 (显示前6位和后4位,中间用*代替)
func Mask(idNumber string) string {
	cleaned := cleanSpaces(idNumber)
	length := len(cleaned)

	if length < 10 {
		return cleaned
	}

	return cleaned[:6] + "********" + cleaned[length-4:]
}

// Hide 隐藏出生日期 (显示XXXX-XX-XX格式)
func Hide(idNumber string) string {
	cleaned := cleanSpaces(idNumber)

	if chineseIDRegex.MatchString(cleaned) {
		// 中国公民身份证: 隐藏月日
		return cleaned[:6] + "XXXXXX" + cleaned[14:]
	}

	if foreignerIDRegex.MatchString(cleaned) {
		// 外国人居留身份证: 隐藏生日部分 (YYMMDD在位置8-11, 4位)
		return cleaned[:8] + "XXXX" + cleaned[12:]
	}

	return cleaned
}

// GetCountryCode 获取国家代码
func GetCountryCode(idNumber string) string {
	cleaned := cleanSpaces(idNumber)
	if len(cleaned) >= 3 && cleaned[0] == '8' {
		return cleaned[1:3]
	}
	return ""
}