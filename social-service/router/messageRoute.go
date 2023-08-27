package service

import (
	"douyin/social-service/middleware"
	"github.com/gin-gonic/gin"
)

func InitMessageRoute(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/message")
	{
		router.POST("/chat", middleware.AuthMiddleware(), middleware.CORSMiddleware())
		router.GET("/action", middleware.AuthMiddleware(), middleware.CORSMiddleware())
	}
	return r
}
