package main

import (
	"douyin/douyin-api/globalinit"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	//导包
	_ "net/http/pprof"
)
//配置
func startPProfServer() {
	go func() {
		zap.S().Info("启动 pprof 服务")
		if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
			zap.S().Errorf("启动 pprof 服务失败: %s", err.Error())
		}
	}()
}
func main() {
	port := 8888
	//初始化logger
	globalinit.InitLogger()
	//初始化routers
	routers := globalinit.Routers()

	//配置pprof性能监控，在本地的6060端口监控
	startPProfServer()

	zap.S().Infof("启动服务器,端口: %d", port)
	if err := routers.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败 ", err.Error())
	}
}
