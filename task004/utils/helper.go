package utils

import (
	"github.com/go-playground/validator/v10"
)

// GetValidationErrors 获取验证错误信息
func GetValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		errs[fieldErr.Field()] = fieldErr.Tag()
	}
	return errs
}

// TruncateString 截断字符串
func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	// 将字符串转换为rune切片以正确处理多字节字符
	runes := []rune(str)
	if len(runes) <= length {
		return str
	}

	return string(runes[:length]) + "..."
}
