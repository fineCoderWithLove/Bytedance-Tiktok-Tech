package router

import (
	"douyin/douyin-api/api"
	"douyin/douyin-api/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
	user全局的路由信息
*/
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/douyin")

	// 添加CORS中间件
	UserRouter.Use(cors.Default())
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("/user/", api.GetUserDetail) //用户详情
		UserRouter.POST("/user/register/", api.UserRegister)  //注册
		UserRouter.POST("/user/login/", api.UserLogin) //登录

		//UserRouter.POST("/publish/action/", api.VideoPublish) //用户投稿视频
		UserRouter.POST("/publish/action/", api.VideoPublish) //测试demo
		UserRouter.GET("/feed/",util.AuthMiddleware(), api.VideoStream)  //视频流信息
		UserRouter.GET("/publish/list/", api.VideoList) //获取用户发布的视频
	}
}