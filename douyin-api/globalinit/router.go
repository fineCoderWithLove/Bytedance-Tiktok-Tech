package globalinit

import (
	"douyin/douyin-api/router"
	"github.com/gin-gonic/gin"

)

/*
	路由的初始化
 */
func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("/v1") //这是一个全局的
 	router.InitUserRouter(ApiGroup)
	return Router
}