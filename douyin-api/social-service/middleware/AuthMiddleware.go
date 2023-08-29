package middleware

import (
	"douyin/douyin-api/social-service/proto"
	"douyin/douyin-api/social-service/utils/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, proto.ErrorMessage{
				Response: proto.Response{
					StatusMsg:  "unauthorized",
					StatusCode: 1,
				},
			})
			return
		}

		_, e := c.Get("userId")
		if e {
			c.Next()
			return
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			log.Printf("AuthMiddleware|token解析错误|%v", err)
			c.JSON(http.StatusUnauthorized, proto.ErrorMessage{
				Response: proto.Response{
					StatusMsg:  "unauthorized",
					StatusCode: 1,
				},
			})
			return
		}
		userID := claims.UserID
		c.Set("userId", userID)
		c.Next()
		return
	}
}
