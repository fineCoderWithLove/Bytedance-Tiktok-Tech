package service

import (
	"douyin/social-service/middleware"
	"github.com/gin-gonic/gin"
)

func InitRelationRoute(r *gin.RouterGroup) gin.IRoutes {

	router := r.Group("/relation")
	{
		router.POST("/action", middleware.AuthMiddleware(), middleware.CORSMiddleware()) // 身份验证+全局跨域
		router.GET("/follow/list", middleware.AuthMiddleware(), middleware.CORSMiddleware())
		router.GET("/follower/list", middleware.AuthMiddleware(), middleware.CORSMiddleware())
		router.GET("/friend/list", middleware.AuthMiddleware(), middleware.CORSMiddleware())
	}
	return r
}
