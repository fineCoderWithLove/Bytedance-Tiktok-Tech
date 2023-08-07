package router

import (
	"douyin/douyin-api/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
	user全局的路由信息
*/
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("/douyin/user/", api.GetUserDetail) // 获取用户的详情信息
		UserRouter.POST("/douyin/user/register/", api.UserRegister) //用户注册
		UserRouter.POST("/douyin/user/login/", api.UserLogin) //用户登录
	}
}
