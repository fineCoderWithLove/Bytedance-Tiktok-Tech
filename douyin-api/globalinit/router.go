package globalinit

import (
	"demotest/douyin-api/router"
	"github.com/gin-gonic/gin"
)

/*
	路由的初始化
*/
func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("") //这是一个全局的
	router.InitFavoriteRouter(ApiGroup)
	router.InitCommentRouter(ApiGroup)
	router.InitUserRouter(ApiGroup)
	return Router
}
