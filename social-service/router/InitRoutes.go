package service

import (
	"douyin/social-service/dal/config"
	"douyin/social-service/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	engine := gin.Default()

	// 配置全局跨域 中间件
	engine.Use(middleware.CORSMiddleware())

	// 路由分组 添加前缀
	group := engine.Group(config.Conf.System.UrlPathPrefix)

	InitMessageRoute(group)
	InitRelationRoute(group)

	return engine
}
