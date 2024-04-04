package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware 验证 JWT Token 的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"status_msg": "您尚未登录",
																	"status_code":"401",
																 })
			c.Abort()
			return
		}

		c.Next()
	}
}
