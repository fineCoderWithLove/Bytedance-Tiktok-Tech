package main

import (
	"douyin/douyin-api/globalinit"
	"fmt"
	"go.uber.org/zap"
)

func main()  {
	port := 8888
	//初始化logger
	globalinit.InitLogger()
	//初始化routers
	routers := globalinit.Routers()
	zap.S().Infof("启动服务器,端口: %d",port)
	if err := routers.Run(fmt.Sprintf(":%d",port)); err!=nil {
		zap.S().Panic("启动失败 ",err.Error())
	}
}
