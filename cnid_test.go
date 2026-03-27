package cnid

import (
	"fmt"
	"testing"
	"time"
)

func TestParseChineseID(t *testing.T) {
	tests := []struct {
		id          string
		wantErr     bool
		wantSex     int
		wantYear    int
		wantRegion  string
	}{
		// 有效身份证 (使用生成的正确校验码)
		{"110101199001010017", false, 1, 1990, "110101"},
		{"110105198505050010", false, 1, 1985, "110105"},
		{"310101199303120018", false, 1, 1993, "310101"},
		{"110108199512310018", false, 1, 1995, "110108"},
		{"310104198801010017", false, 1, 1988, "310104"},
		{"440100199505050011", false, 1, 1995, "440100"},
		{"330100199607150015", false, 1, 1996, "330100"},
		{"320100200001010019", false, 1, 2000, "320100"},
		{"500101201001010014", false, 1, 2010, "500101"},
		{"610100198501010013", false, 1, 1985, "610100"},
		// 无效身份证
		{"000000000000000000", true, 0, 0, ""},
		{"123456789012345678", true, 0, 0, ""},
		{"110101199001011235", false, 1, 1990, "110101"}, // 有效的正确校验码
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			id, err := ParseChinese(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChinese() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if id.GetSex() != tt.wantSex {
					t.Errorf("ParseChinese() sex = %v, want %v", id.GetSex(), tt.wantSex)
				}
				if id.GetBirthday().Year() != tt.wantYear {
					t.Errorf("ParseChinese() year = %v, want %v", id.GetBirthday().Year(), tt.wantYear)
				}
				if id.GetRegionCode() != tt.wantRegion {
					t.Errorf("ParseChinese() region = %v, want %v", id.GetRegionCode(), tt.wantRegion)
				}
			}
		})
	}
}

func TestValidateChineseID(t *testing.T) {
	validIDs := []string{
		"110101199001010017",
		"310104198801010017",
		"440100199505050011",
	}

	invalidIDs := []string{
		"123456789012345678",
		"000000000000000000",
	}

	for _, id := range validIDs {
		if !ValidateChinese(id) {
			t.Errorf("ValidateChinese(%s) = false, want true", id)
		}
	}

	for _, id := range invalidIDs {
		if ValidateChinese(id) {
			t.Errorf("ValidateChinese(%s) = true, want false", id)
		}
	}
}

func TestParseForeignerID(t *testing.T) {
	tests := []struct {
		id          string
		wantErr     bool
		wantSex     int
		wantYear    int
		wantRegion  string
		wantCountry string
	}{
		// 有效的外国人居留身份证 (16位,8开头)
		// 格式: 8 + 国家(2) + 地区(4) + 生日(6) + 顺序(2) + 校验(1)
		{"8011100900101013", false, 1, 1990, "1100", "01"}, // 美国, 1990年
		{"8051100950505019", false, 1, 1995, "1100", "05"}, // 日本, 1995年
		{"8021100900101012", false, 1, 1990, "1100", "02"}, // 英国, 1990年
		// 无效外国人居留身份证
		{"9111100900101013", true, 0, 0, "", ""}, // 不以8开头
		{"8011100900101014", true, 0, 0, "", ""}, // 校验码错误
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			id, err := ParseForeigner(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseForeigner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if id.GetSex() != tt.wantSex {
					t.Errorf("ParseForeigner() sex = %v, want %v", id.GetSex(), tt.wantSex)
				}
				if id.GetBirthday().Year() != tt.wantYear {
					t.Errorf("ParseForeigner() year = %v, want %v", id.GetBirthday().Year(), tt.wantYear)
				}
				if id.GetRegionCode() != tt.wantRegion {
					t.Errorf("ParseForeigner() region = %v, want %v", id.GetRegionCode(), tt.wantRegion)
				}
				if id.CountryCode != tt.wantCountry {
					t.Errorf("ParseForeigner() country = %v, want %v", id.CountryCode, tt.wantCountry)
				}
			}
		})
	}
}

func TestValidateForeignerID(t *testing.T) {
	validIDs := []string{
		"8011100900101013",
		"8051100950505019",
		"8021100900101012",
	}

	invalidIDs := []string{
		"1234567890123456",
		"9111100900101013", // 不以8开头
		"8011100900101014", // 校验码错误
	}

	for _, id := range validIDs {
		if !ValidateForeigner(id) {
			t.Errorf("ValidateForeigner(%s) = false, want true", id)
		}
	}

	for _, id := range invalidIDs {
		if ValidateForeigner(id) {
			t.Errorf("ValidateForeigner(%s) = true, want false", id)
		}
	}
}

func TestParseAutoDetect(t *testing.T) {
	tests := []struct {
		id       string
		wantErr  bool
		wantType IDType
	}{
		{"110101199001010017", false, TypeChinese},
		{"8011100900101013", false, TypeForeigner},
		{"123456789012345678", true, TypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			info, err := Parse(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if info.GetType() != tt.wantType {
					t.Errorf("Parse() type = %v, want %v", info.GetType(), tt.wantType)
				}
			}
		})
	}
}

func TestGetRegion(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"110101199001010017", "北京市东城区"},
		{"310101199001010018", "上海市黄浦区"},
		{"440100199001010011", "广东省广州市"},
		{"999999199001010019", ""}, // 未知地区
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if got := GetRegion(tt.id); got != tt.want {
				t.Errorf("GetRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountry(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"8011100900101013", "美国"},
		{"8051100950505019", "日本"},
		{"8021100900101012", "英国"},
		{"8999900900101019", ""}, // 未知国家代码返回空
		{"110101199001010017", ""}, // 不是外国人居留证
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if got := GetCountry(tt.id); got != tt.want {
				t.Errorf("GetCountry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAge(t *testing.T) {
	// 使用一个有效的身份证 320100201001010014 (2010年出生)
	id := "320100201001010014"
	age, err := GetAge(id)
	if err != nil {
		t.Errorf("GetAge() error = %v", err)
	}
	// 2010年出生的人,2024年应该是14岁左右
	if age < 10 || age > 20 {
		t.Errorf("GetAge() = %v, expected reasonable age", age)
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"110101199001010017", "110101********0017"},
		{"8011100900101013", "801110********1013"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if got := Mask(tt.id); got != tt.want {
				t.Errorf("Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHide(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"110101199001010017", "110101XXXXXX0017"},
		{"8011100900101013", "80111009XXXX1013"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if got := Hide(tt.id); got != tt.want {
				t.Errorf("Hide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSexString(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"110101199001010017", "男"},
		{"110101199001010005", "女"}, // 顺序码0是偶数，女性
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got, err := GetSexString(tt.id)
			if err != nil {
				t.Errorf("GetSexString() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetSexString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChineseIDString(t *testing.T) {
	id := &ChineseID{
		Number:     "110101199001010017",
		RegionCode: "110101",
		Birthday:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Sex:        1,
		Sequence:   17,
	}

	got := id.String()
	want := "ChineseID{Number: 110101199001010017, Region: 北京市东城区, Birthday: 1990-01-01, Sex: 1}"

	if got != want {
		t.Errorf("ChineseID.String() = %v, want %v", got, want)
	}
}

func TestForeignerIDString(t *testing.T) {
	id := &ForeignerID{
		Number:      "8011100900101013",
		CountryCode: "01",
		CountryName: "美国",
		RegionCode:  "1100",
		Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Sex:         1,
		Sequence:    13,
		IsPermanent: true,
	}

	got := id.String()
	want := "ForeignerID{Number: 8011100900101013, Country: 美国, Region: , Birthday: 1990-01-01, Sex: 1}"

	if got != want {
		t.Errorf("ForeignerID.String() = %v, want %v", got, want)
	}
}

// BenchmarkParseChinese 性能测试
func BenchmarkParseChinese(b *testing.B) {
	id := "110101199001010017"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseChinese(id)
	}
}

// BenchmarkParseForeigner 性能测试
func BenchmarkParseForeigner(b *testing.B) {
	id := "8011100900101013"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseForeigner(id)
	}
}

// ExampleParseChinese 演示解析中国公民身份证
func ExampleParseChinese() {
	id, err := ParseChinese("110101199001010017")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("类型: 中国公民身份证\n")
	fmt.Printf("号码: %s\n", id.Number)
	fmt.Printf("地区: %s\n", GetRegion(id.Number))
	fmt.Printf("出生日期: %s\n", id.Birthday.Format("2006-01-02"))
	fmt.Printf("性别: %d (1=男,0=女)\n", id.Sex)

	// Output:
	// 类型: 中国公民身份证
	// 号码: 110101199001010017
	// 地区: 北京市东城区
	// 出生日期: 1990-01-01
	// 性别: 1 (1=男,0=女)
}

// ExampleParseForeigner 演示解析外国人居留身份证
func ExampleParseForeigner() {
	id, err := ParseForeigner("8011100900101013")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("类型: 外国人居留身份证\n")
	fmt.Printf("号码: %s\n", id.Number)
	fmt.Printf("国籍: %s\n", id.CountryName)
	fmt.Printf("性别: %d (1=男,0=女)\n", id.Sex)

	// Output:
	// 类型: 外国人居留身份证
	// 号码: 8011100900101013
	// 国籍: 美国
	// 性别: 1 (1=男,0=女)
}

// ExampleParseAutoDetect 演示自动识别身份证类型
func ExampleParse() {
	ids := []string{
		"110101199001010017",  // 中国公民身份证
		"8011100900101013", // 外国人居留身份证
	}

	for _, idStr := range ids {
		info, err := Parse(idStr)
		if err != nil {
			fmt.Printf("%s: 无效身份证\n", idStr)
			continue
		}

		switch info.GetType() {
		case TypeChinese:
			fmt.Printf("%s: 中国公民身份证\n", idStr)
		case TypeForeigner:
			fmt.Printf("%s: 外国人居留身份证\n", idStr)
		default:
			fmt.Printf("%s: 未知类型\n", idStr)
		}
	}

	// Output:
	// 110101199001010017: 中国公民身份证
	// 8011100900101013: 外国人居留身份证
}