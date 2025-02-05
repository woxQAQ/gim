package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt算法对密码进行哈希处理
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword 比较密码是否匹配
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePassword 验证密码是否正确
func ValidatePassword(password, hashedPassword string) bool {
	err := ComparePassword(hashedPassword, password)
	return err == nil
}
