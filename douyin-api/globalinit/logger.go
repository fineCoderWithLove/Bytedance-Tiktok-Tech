package globalinit

import "go.uber.org/zap"

<<<<<<< HEAD
func InitLogger()  {
	//初始化全局日志信息
	logger,_ := zap.NewDevelopment()
=======
func InitLogger() {
	//初始化全局日志信息
	logger, _ := zap.NewDevelopment()
>>>>>>> cba9c25843da297a4159b839c47e609847fe7bed
	zap.ReplaceGlobals(logger)
}
