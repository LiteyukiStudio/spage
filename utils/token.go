package utils

import (
	"crypto/rand"
	"github.com/LiteyukiStudio/spage/config"
	"github.com/LiteyukiStudio/spage/spage/models"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type TokenType struct{}

var Token = TokenType{}

type Claims struct {
	jwt.RegisteredClaims
	UserID   uint `json:"user_id"`  // 用户ID，用于身份验证
	TokenID  uint `json:"token_id"` // 令牌ID，用于服务端会话维持
	Stateful bool `json:"stateful"` // 是否为有状态Token
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[index.Int64()]
	}
	return string(result)
}

// CreateJsonWebToken 生成用户会话令牌（默认24小时有效）
func (TokenType) CreateJsonWebToken(userID uint, duration time.Duration, stateful bool, persistentHandler func(uint) (*models.JsonWebToken, error)) (string, error) {
	var tokenModel *models.JsonWebToken
	var err error
	if stateful {
		tokenModel, err = persistentHandler(userID)
		if err != nil {
			return "", err
		}
	} else {
		tokenModel = &models.JsonWebToken{
			Model: gorm.Model{
				ID: 0,
			},
		}
	}

	claims := Claims{
		UserID:   userID,
		TokenID:  tokenModel.ID,
		Stateful: stateful,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtSecret))
}

// ParseJsonWebToken 解析JWT令牌
func (TokenType) ParseJsonWebToken(tokenString string, revokeChecker func(uint) bool) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// 有状态token被吊销也视为过期
	if claims.Stateful {
		if revokeChecker(claims.UserID) {
			return nil, jwt.ErrTokenExpired
		}
	}
	return claims, nil
}

// CreateApiToken 生成API令牌
func (TokenType) CreateApiToken(userID uint, duration time.Duration, persistentHandler func(*models.ApiToken) error) (string, error) {
	expiration := time.Now().Add(duration)
	apiToken := &models.ApiToken{
		UserID:    userID,
		Token:     "spat_" + GenerateRandomString(32), // 生成随机字符串作为令牌
		ExpiresAt: expiration,
	}
	err := persistentHandler(apiToken)
	if err != nil {
		return "", err
	}
	return apiToken.Token, nil
}

// ParseApiToken 解析API令牌
func (TokenType) ParseApiToken(tokenString string, isValidFunc func(string) (*models.ApiToken, error)) (*models.ApiToken, error) {
	apiToken, err := isValidFunc(tokenString)
	if err != nil {
		return nil, err
	}
	if apiToken == nil || apiToken.ExpiresAt.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	return apiToken, nil
}
