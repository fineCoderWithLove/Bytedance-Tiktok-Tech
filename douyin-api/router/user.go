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
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("list", api.GetUserVideo)
	}
}
