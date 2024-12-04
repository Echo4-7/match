package v1

import "regexp"

// IsEmail 判断是否是邮箱
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// IsTelNum 判断是否是电话号码
func IsTelNum(telNum string) bool {
	pattern := `^1[3-9]\d{9}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(telNum)
}

// IsUserId 判断是否是 17 位 userId
func IsUserId(userId string) bool {
	pattern := `^\d{17}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(userId)
}
