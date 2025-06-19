package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type PasswordType struct {
}

var Password = PasswordType{}

// HashPassword 密码哈希函数
func (u *PasswordType) HashPassword(password string, salt string) (string, error) {
	saltedPassword := Password.addSalt(password, salt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 验证密码
func (u *PasswordType) VerifyPassword(password, hashedPassword string, salt string) bool {
	if len(hashedPassword) == 0 || len(salt) == 0 {
		// 防止oidc空密码出问题
		return false
	}
	saltedPassword := Password.addSalt(password, salt)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	return err == nil
}

// addSalt 加盐函数
// password: 原始密码
// salt: 盐值
// 返回值: 加盐后的密码
func (u *PasswordType) addSalt(password string, salt string) string {
	combined := password + salt
	hash := sha256.New()
	hash.Write([]byte(combined))
	return hex.EncodeToString(hash.Sum(nil))
}

// CheckPasswordComplexity 根据指定级别检查密码复杂度
// password: 待检查的密码
// level: 复杂度级别(1-4)
// 返回值: 是否满足复杂度要求
func (u *PasswordType) CheckPasswordComplexity(password string, level int) bool {
	if len(password) <= 8 {
		return false
	}

	// 定义各种字符类型的检查标志
	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
		typesUsed  = 0
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char) && !hasLower:
			hasLower = true
			typesUsed++
		case unicode.IsUpper(char) && !hasUpper:
			hasUpper = true
			typesUsed++
		case unicode.IsDigit(char) && !hasDigit:
			hasDigit = true
			typesUsed++
		case unicode.IsPunct(char) || unicode.IsSymbol(char) && !hasSpecial:
			hasSpecial = true
			typesUsed++
		}
	}

	return typesUsed >= level
}
