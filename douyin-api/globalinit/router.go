package globalinit

import (
	"douyin/douyin-api/router"
	"github.com/gin-gonic/gin"
<<<<<<< HEAD

=======
>>>>>>> cba9c25843da297a4159b839c47e609847fe7bed
)

/*
	路由的初始化
<<<<<<< HEAD
 */
func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("") //这是一个全局的
 	router.InitUserRouter(ApiGroup)
	return Router
}
=======
*/
func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("") //这是一个全局的
	router.InitFavoriteRouter(ApiGroup)
	router.InitCommentRouter(ApiGroup)
	router.InitUserRouter(ApiGroup)
	return Router
}
>>>>>>> cba9c25843da297a4159b839c47e609847fe7bed
