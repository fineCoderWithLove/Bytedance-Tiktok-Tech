package global

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var InfoLogger *zap.Logger
var ErrLogger *zap.Logger
var ConsoleLogger *zap.Logger

func init() {
	InitLogger()
	defer InfoLogger.Sync()
	defer ErrLogger.Sync()
	defer ConsoleLogger.Sync()
	ConsoleLogger.Info("日志服务启动")
}

func InitLogger() {
	infoWriteSyncer := getLogWriter("info")
	errWriteSyncer := getLogWriter("error")
	encoder := getEncoder()
	infoCore := zapcore.NewCore(encoder, infoWriteSyncer, zapcore.InfoLevel)
	errCore := zapcore.NewCore(encoder, errWriteSyncer, zapcore.ErrorLevel)

	InfoLogger = zap.New(infoCore, zap.AddCaller())
	ErrLogger = zap.New(errCore, zap.AddCaller())

	//log, _ = zap.NewProduction()//生产环境
	logger, err := zap.NewDevelopment() //开发环境
	if err != nil {
		panic(err)
	}
	ConsoleLogger = logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(infoOrErr string) zapcore.WriteSyncer {
	var lumberJackLogger *lumberjack.Logger
	if infoOrErr == "info" {
		lumberJackLogger = &lumberjack.Logger{
			Filename:   "./log/info/info.log",
			MaxSize:    1,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   false,
			LocalTime:  true,
		}
	} else if infoOrErr == "error" {
		lumberJackLogger = &lumberjack.Logger{
			Filename:   "./log/error/error.log",
			MaxSize:    1,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   false,
			LocalTime:  true,
		}
	}
	return zapcore.AddSync(lumberJackLogger)
}
