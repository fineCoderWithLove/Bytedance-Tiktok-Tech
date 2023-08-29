package service

import (
	"douyin/douyin-api/social-service/handler"
	util "douyin/douyin-api/social-service/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRelationRoute(r *gin.RouterGroup) gin.IRoutes {
	relationController := handler.NewRelationController()
	router := r.Group("/douyin/relation")
	zap.S().Info("配置关注相关的url")
	{
		router.POST("/action", util.AuthMiddleware(), relationController.FollowAction)
		router.GET("/follow/list", util.AuthMiddleware(), relationController.GetFollow)
		router.GET("/follower/list", util.AuthMiddleware(), relationController.GetFollower)
		router.GET("/friend/list", util.AuthMiddleware(), relationController.GetFriendList)
	}
	return r
}
