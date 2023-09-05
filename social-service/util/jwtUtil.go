// util/jwtUtil.go

package util

import (
	"demotest/social-service/global"
	"github.com/dgrijalva/jwt-go"

	"time"
)

const (
	// 在生产环境中，请替换为安全的密钥，不要将密钥硬编码在代码中
	jwtSecret = "your_jwt_secret_key"
)

// GenerateToken 生成 JWT token
func GenerateToken(userId int64, expiration time.Duration) (string, error) {
	// 创建一个新的 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(), // Token 的过期时间（在此示例中为expiration）
		},
	})

	// 使用密钥对 token 进行签名
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	// 将 Token 存储到 Redis 中，设置过期时间为expiration
	err = SetToken(tokenString, expiration)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CustomClaims 自定义 JWT token 中的声明信息
type CustomClaims struct {
	UserId int64 `json:"user_id"` // 用户ID
	jwt.StandardClaims
}

// VerifyToken 验证 JWT token
func VerifyToken(tokenString string) (*CustomClaims, error) {
	// 解析并验证 token 的有效性
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证 token 的声明信息
	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// GetToken 从 Redis 中获取 Token
func GetToken(token string) (string, error) {
	return global.RS.Get(token).Result()
}

// SetToken 将 Token 存储到 Redis 中
func SetToken(token string, expiration time.Duration) error {
	return global.RS.Set(token, token, expiration).Err()
}
