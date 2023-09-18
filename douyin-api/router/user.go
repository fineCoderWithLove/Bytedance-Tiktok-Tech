package router

import (
	"demotest/douyin-api/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"demotest/douyin-api/util"
)

/*
	user全局的路由信息
 	需要敏感词过滤的中间件
	1.用户注册
 	2.用户发布视频时候的标题
	3.用户发送消息的行为
*/
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/douyin")

	// 添加CORS中间件，配置路由信息，base-service and social-service路由信息
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("/user/", api.GetUserDetail) //用户详情
		UserRouter.POST("/user/register/",util.UserRegisterMiddleware(),api.UserRegister)  //注册
		UserRouter.POST("/user/login/", api.UserLogin) //登录
//util.AuthMiddleware(),

		UserRouter.POST("/publish/action/",util.PublishVideoMiddleware(),api.VideoPublish) //发布视频
		UserRouter.GET("/feed/", api.VideoStream)  //视频流信息
		UserRouter.GET("/publish/list/", api.VideoList) //获取用户发布的视频

		UserRouter.POST("/relation/action/", api.UserRelationAction) //用户关注的行为
		UserRouter.GET("/relation/follow/list/", api.UserRelationFollowList) //获取用户关注的列表
		UserRouter.GET("/relation/follower/list/", api.UserRelationFollowerList) //获取用户的粉丝列表
		UserRouter.GET("/relation/friend/list/", api.UserRelationFriendList) //获取用户的好友列表

		UserRouter.POST("/message/action/",util.MessageActionMiddleware(), api.UserMessageAction) //用户发送消息的行为
		UserRouter.GET("/message/chat/", api.UserMessageChat) //获取用户的聊天记录

	}
}