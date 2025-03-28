package middleware

import (
	"chat/pkg/auth"
	myerrors "chat/pkg/errors"
	"chat/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// jwt中间件
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的Authorization字段
		authHeader := c.GetHeader("Authorization")
		// 如果Authorization字段为空，返回401错误信息，并终止请求处理
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Fail(myerrors.ErrCodeMissingToken, "未携带token"))
			c.Abort()
			return
		}

		// 从Authorization字段中提取token字符串，移除"Bearer "前缀
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			if errors.Is(myerrors.ErrTokenExpired, err) {
				c.JSON(http.StatusUnauthorized, response.Fail(myerrors.ErrCodeTokenExpired, "token已失效"))
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, response.Fail(myerrors.ErrCodeInvalidToken, "token错误"))
			c.Abort()
			return
		}

		// 将解析出的用户名设置到请求上下文中，供后续处理函数使用
		c.Set("username", claims.Username)
		c.Set("userid", claims.UserId)
		c.Set("token", tokenString)

		c.Next()

	}

}
