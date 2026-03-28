package cnid

import (
	"testing"
	"time"
)

// 测试中国居民身份证新版（18 位）
func TestValidateResidentNew18(t *testing.T) {
	// 使用生成的有效身份证进行测试
	validID := GenerateResident("110105", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "男")

	tests := []struct {
		name     string
		idNumber string
		want     bool
	}{
		{"有效身份证 1", "11010519491231002X", true},
		{"有效身份证 2", validID, true},
		{"有效身份证 3", "440300199001011234", false}, // 校验码可能不匹配
		{"无效格式", "11010519491231002A", false},
		{"长度错误", "11010519491231002", false},
		{"空字符串", "", false},
		{"地区代码无效", "00010519491231002X", false},
		{"日期无效", "11010520230230002X", false}, // 2 月 30 日不存在
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.idNumber); got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.idNumber, got, tt.want)
			}
		})
	}
}

// 测试中国居民身份证旧版（15 位）
func TestValidateResidentOld15(t *testing.T) {
	tests := []struct {
		name     string
		idNumber string
		want     bool
	}{
		{"有效身份证 1", "110105491231002", true},
		{"有效身份证 2", "440300900101123", true},
		{"无效格式", "11010549123100A", false},
		{"长度错误", "11010549123100", false},
		{"地区代码无效", "000105491231002", false},
		{"日期无效", "110105230230002", false}, // 2 月 30 日不存在
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.idNumber); got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.idNumber, got, tt.want)
			}
		})
	}
}

// 测试外国人永久居留身份证新版（18 位）
func TestValidateForeignNew18(t *testing.T) {
	// 使用生成的有效身份证进行测试
	validID1 := GenerateForeignNew("USA", "11", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "男")
	validID2 := GenerateForeignNew("CHN", "31", time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC), "女")

	tests := []struct {
		name     string
		idNumber string
		want     bool
	}{
		{"有效永居证 1", validID1, true},
		{"有效永居证 2", validID2, true},
		{"不以 9 开头", "8USA11200001011231", false},
		{"长度错误", "9USA1120000101123", false},
		{"格式错误", "9US111200001011231", false},
		{"校验码错误", "9USA1119900101123X", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.idNumber); got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.idNumber, got, tt.want)
			}
		})
	}
}

// 测试外国人永久居留身份证旧版（15 位）
func TestValidateForeignOld15(t *testing.T) {
	tests := []struct {
		name     string
		idNumber string
		want     bool
	}{
		{"有效永居证 1", "USA199001011234", true},
		{"有效永居证 2", "CHN200001011234", true},
		{"格式错误", "1SA199001011234", false}, // 首字符必须是字母
		{"长度错误", "USA19900101123", false},
		{"小写字母", "usa199001011234", true}, // 小写会被转换为大写验证
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.idNumber); got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.idNumber, got, tt.want)
			}
		})
	}
}

// 测试 GetType 函数
func TestGetType(t *testing.T) {
	validForeignID := GenerateForeignNew("USA", "11", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "男")

	tests := []struct {
		name     string
		idNumber string
		wantType int
	}{
		{"中国居民身份证新版", "11010519491231002X", TypeResidentNew18},
		{"中国居民身份证旧版", "110105491231002", TypeResidentOld15},
		{"外国人永居证新版", validForeignID, TypeForeignNew18},
		{"外国人永居证旧版", "USA199001011234", TypeForeignOld15},
		{"未知类型", "invalid", TypeUnknown},
		{"空字符串", "", TypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetType(tt.idNumber); got != tt.wantType {
				t.Errorf("GetType(%q) = %v, want %v", tt.idNumber, got, tt.wantType)
			}
		})
	}
}

// 测试 Parse 函数 - 中国居民身份证新版
func TestParseResidentNew18(t *testing.T) {
	// 生成一个有效的测试身份证
	idNumber := GenerateResident("110105", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "男")

	info, err := Parse(idNumber)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", idNumber, err)
	}

	if info.Type != TypeResidentNew18 {
		t.Errorf("Type = %v, want %v", info.Type, TypeResidentNew18)
	}

	if info.IDNumber != idNumber {
		t.Errorf("IDNumber = %v, want %v", info.IDNumber, idNumber)
	}

	if info.AreaCode != "110105" {
		t.Errorf("AreaCode = %v, want 110105", info.AreaCode)
	}

	if info.BirthDate.Year() != 1990 || info.BirthDate.Month() != 1 || info.BirthDate.Day() != 1 {
		t.Errorf("BirthDate = %v, want 1990-01-01", info.BirthDate)
	}

	if info.Gender != "男" {
		t.Errorf("Gender = %v, want 男", info.Gender)
	}
}

// 测试 Parse 函数 - 中国居民身份证旧版
func TestParseResidentOld15(t *testing.T) {
	idNumber := "110105900101123"

	info, err := Parse(idNumber)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", idNumber, err)
	}

	if info.Type != TypeResidentOld15 {
		t.Errorf("Type = %v, want %v", info.Type, TypeResidentOld15)
	}

	if info.AreaCode != "110105" {
		t.Errorf("AreaCode = %v, want 110105", info.AreaCode)
	}

	// 旧版身份证年份默认为 19xx
	if info.BirthDate.Year() != 1990 || info.BirthDate.Month() != 1 || info.BirthDate.Day() != 1 {
		t.Errorf("BirthDate = %v, want 1990-01-01", info.BirthDate)
	}

	// 第 15 位是 3（奇数），应该是男性
	if info.Gender != "男" {
		t.Errorf("Gender = %v, want 男", info.Gender)
	}
}

// 测试 Parse 函数 - 外国人永久居留身份证新版
func TestParseForeignNew18(t *testing.T) {
	// 生成一个有效的测试永居证
	idNumber := GenerateForeignNew("USA", "11", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "女")

	if !Validate(idNumber) {
		t.Fatalf("Generated ID %q is not valid", idNumber)
	}

	info, err := Parse(idNumber)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", idNumber, err)
	}

	if info.Type != TypeForeignNew18 {
		t.Errorf("Type = %v, want %v", info.Type, TypeForeignNew18)
	}

	if info.Nationality != "USA" {
		t.Errorf("Nationality = %v, want USA", info.Nationality)
	}

	if info.IssuePlace != "11" {
		t.Errorf("IssuePlace = %v, want 11", info.IssuePlace)
	}

	if info.BirthDate.Year() != 1990 || info.BirthDate.Month() != 1 || info.BirthDate.Day() != 1 {
		t.Errorf("BirthDate = %v, want 1990-01-01", info.BirthDate)
	}

	if info.Gender != "女" {
		t.Errorf("Gender = %v, want 女", info.Gender)
	}
}

// 测试 Parse 函数 - 外国人永久居留身份证旧版
func TestParseForeignOld15(t *testing.T) {
	idNumber := "USA199001011234"

	info, err := Parse(idNumber)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", idNumber, err)
	}

	if info.Type != TypeForeignOld15 {
		t.Errorf("Type = %v, want %v", info.Type, TypeForeignOld15)
	}

	if info.Nationality != "USA" {
		t.Errorf("Nationality = %v, want USA", info.Nationality)
	}

	// 旧版只有年月，日默认为 01
	if info.BirthDate.Year() != 1990 || info.BirthDate.Month() != 1 || info.BirthDate.Day() != 1 {
		t.Errorf("BirthDate = %v, want 1990-01-01", info.BirthDate)
	}
}

// 测试 UpgradeOld15To18 函数
func TestUpgradeOld15To18(t *testing.T) {
	tests := []struct {
		name   string
		id15   string
		want18 string
	}{
		{"升级 1", "110105491231002", "11010519491231002X"},
		{"升级 2", "110105900101123", "110105199001011232"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpgradeOld15To18(tt.id15)
			if err != nil {
				t.Fatalf("UpgradeOld15To18(%q) error: %v", tt.id15, err)
			}
			if got != tt.want18 {
				t.Errorf("UpgradeOld15To18(%q) = %v, want %v", tt.id15, got, tt.want18)
			}

			// 验证升级后的身份证是否有效
			if !Validate(got) {
				t.Errorf("Upgraded ID %q is not valid", got)
			}
		})
	}
}

// 测试 GenerateResident 函数
func TestGenerateResident(t *testing.T) {
	// 测试完全随机生成
	id1 := GenerateResident("", time.Time{}, "")
	if !Validate(id1) {
		t.Errorf("Generated ID %q is not valid", id1)
	}
	if GetType(id1) != TypeResidentNew18 {
		t.Errorf("Generated ID type = %v, want %v", GetType(id1), TypeResidentNew18)
	}

	// 测试指定参数生成
	birthDate := time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC)
	id2 := GenerateResident("440300", birthDate, "女")
	if !Validate(id2) {
		t.Errorf("Generated ID %q is not valid", id2)
	}

	info, _ := Parse(id2)
	if info.AreaCode != "440300" {
		t.Errorf("AreaCode = %v, want 440300", info.AreaCode)
	}
	if info.Gender != "女" {
		t.Errorf("Gender = %v, want 女", info.Gender)
	}
}

// 测试 GenerateForeignNew 函数
func TestGenerateForeignNew(t *testing.T) {
	// 测试完全随机生成
	id1 := GenerateForeignNew("", "", time.Time{}, "")
	if !Validate(id1) {
		t.Errorf("Generated ID %q is not valid", id1)
	}
	if GetType(id1) != TypeForeignNew18 {
		t.Errorf("Generated ID type = %v, want %v", GetType(id1), TypeForeignNew18)
	}

	// 测试指定参数生成
	birthDate := time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC)
	id2 := GenerateForeignNew("CHN", "31", birthDate, "男")
	if !Validate(id2) {
		t.Errorf("Generated ID %q is not valid", id2)
	}

	info, _ := Parse(id2)
	if info.Nationality != "CHN" {
		t.Errorf("Nationality = %v, want CHN", info.Nationality)
	}
	if info.Gender != "男" {
		t.Errorf("Gender = %v, want 男", info.Gender)
	}
}

// 测试校验码计算
func TestCalculateCheckCode(t *testing.T) {
	tests := []struct {
		body       string
		wantCheck  string
	}{
		{"11010519491231002", "X"},
		{"11010519900101123", calculateCheckCode("11010519900101123")},
		{"911000019900101123", calculateCheckCode("911000019900101123")},
	}

	for _, tt := range tests {
		got := calculateCheckCode(tt.body)
		if got != tt.wantCheck {
			t.Errorf("calculateCheckCode(%q) = %v, want %v", tt.body, got, tt.wantCheck)
		}
	}
}

// 测试 GetTypeName 函数
func TestGetTypeName(t *testing.T) {
	tests := []struct {
		idType int
		want   string
	}{
		{TypeResidentOld15, "中国居民身份证旧版（15 位）"},
		{TypeResidentNew18, "中国居民身份证新版（18 位）"},
		{TypeForeignOld15, "外国人永久居留身份证旧版（15 位）"},
		{TypeForeignNew18, "外国人永久居留身份证新版（18 位）"},
		{TypeUnknown, "未知类型"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := GetTypeName(tt.idType); got != tt.want {
				t.Errorf("GetTypeName(%v) = %v, want %v", tt.idType, got, tt.want)
			}
		})
	}
}

// 测试 IDInfo String 方法
func TestIDInfoString(t *testing.T) {
	idNumber := GenerateResident("110105", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), "男")
	info, _ := Parse(idNumber)

	str := info.String()
	if str == "" {
		t.Error("IDInfo.String() should not return empty string")
	}

	// 检查是否包含关键字段
	expectedFields := []string{"身份证类型", "身份证号码", "出生日期", "性别", "地区代码"}
	for _, field := range expectedFields {
		if !contains(str, field) {
			t.Errorf("IDInfo.String() should contain %q", field)
		}
	}
}

// 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// 测试大小写不敏感
func TestCaseInsensitive(t *testing.T) {
	// 小写 x 应该被接受
	idLower := "11010519491231002x"
	idUpper := "11010519491231002X"

	if !Validate(idLower) {
		t.Errorf("Validate(%q) should be true (case insensitive)", idLower)
	}

	if !Validate(idUpper) {
		t.Errorf("Validate(%q) should be true", idUpper)
	}

	// 解析后应该都转换为大写
	infoLower, _ := Parse(idLower)
	infoUpper, _ := Parse(idUpper)

	if infoLower.IDNumber != infoUpper.IDNumber {
		t.Errorf("Parsed ID numbers should be equal after normalization")
	}
}

// 测试空白字符处理
func TestWhitespaceHandling(t *testing.T) {
	idWithSpace := " 11010519491231002X "
	idClean := "11010519491231002X"

	if !Validate(idWithSpace) {
		t.Errorf("Validate(%q) should handle whitespace", idWithSpace)
	}

	infoWithSpace, _ := Parse(idWithSpace)
	infoClean, _ := Parse(idClean)

	if infoWithSpace.IDNumber != infoClean.IDNumber {
		t.Errorf("Parsed ID numbers should be equal after trimming whitespace")
	}
}
