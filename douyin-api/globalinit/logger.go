package globalinit

import "go.uber.org/zap"

func InitLogger() {
	//初始化全局日志信息
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
