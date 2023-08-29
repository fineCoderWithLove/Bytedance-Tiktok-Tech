package service

import (
	"douyin/douyin-api/social-service/handler"
	util "douyin/douyin-api/social-service/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitMessageRoute(r *gin.RouterGroup) gin.IRoutes {
	chatController := handler.NewChatController()
	router := r.Group("/douyin/message")
	zap.S().Info("配置消息相关的url")
	{
		router.POST("/chat", util.AuthMiddleware(), chatController.ChatPost)
		router.GET("/action", util.AuthMiddleware(), chatController.ListChatMessage)
	}
	return r
}
