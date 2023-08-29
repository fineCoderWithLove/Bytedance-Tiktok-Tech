package globalinit

import (
	"douyin/douyin-api/router"
	router2 "douyin/douyin-api/social-service/router"
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
	router2.InitMessageRoute(ApiGroup)
	router2.InitRelationRoute(ApiGroup)
	
	return Router
}
