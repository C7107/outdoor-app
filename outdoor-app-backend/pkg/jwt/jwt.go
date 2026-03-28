package jwt

import (
	"errors"
	"outdoor-app-backend/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ⚠️ 实际开发中这个秘钥应该放在 .env 里，后面再改
// 不在这里初始化，改为函数内获取
func GetJWTSecret() []byte {
	if configs.AppConfig == nil {
		panic("配置未初始化，请先调用 configs.InitConfig()")
	}
	return []byte(configs.AppConfig.JWT.Secret)
}

// CustomClaims 自定义载荷，保存用户的核心信息
type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, email string) (string, error) {
	// 设置 Token 的过期时间为 7 天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "outdoor_app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret())
}

// ParseToken 解析并校验 JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
