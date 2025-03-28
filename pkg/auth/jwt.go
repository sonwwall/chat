package auth

import (
	"chat/internal/global"
	myerrors "chat/pkg/errors"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// jwt
var jwtkey = []byte("sonwwall")

type Claims struct {
	Username string `json:"username"`
	UserId   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成token
func GenerateToken(username string, userid uint) (string, error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(time.Hour * 24)) //设置过期时间
	claims := &Claims{
		Username: username,
		UserId:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// 验证token

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	// 尝试使用提供的密钥解析JWT令牌字符串，并验证其签名。
	// 如果解析成功且签名有效，claims参数将填充解析后的声明
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})

	// 如果解析或验证过程中发生错误，返回错误。
	if err != nil {
		return nil, err
	}
	// 如果令牌被成功解析但无效（例如，签名无效），返回错误。
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// 检查 JWT 是否在黑名单中
	ctx := context.Background()
	blacklisted, err := global.Redis.Get(ctx, tokenString).Result()
	if err == nil && blacklisted == "true" {
		return nil, myerrors.ErrTokenExpired
	}

	return claims, nil

}

// 将jwt添加到黑名单
func AddTokenToBlacklist(tokenString string) error {
	ctx := context.Background()
	expiration := time.Hour * 24 //设置过期时间
	// 使用 global.Redis.Set 方法将 JWT 令牌添加到 Redis 数据库中
	// 第三个参数 "true" 表示在 Redis 中存储的值，这里简单地使用字符串 "true" 来标记该令牌已被列入黑名单。
	err := global.Redis.Set(ctx, tokenString, "true", expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
