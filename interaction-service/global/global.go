package global

import (
	"demotest/douyin-api/globalinit/constant"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
)
var (
	RS *redis.Client
)

func init() {
	/*
		mysql的全局连接
	*/
	dsn := constant.MySQL_DSN
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for log
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	/*
		redis的全局连接
	*/
	RS = redis.NewClient(&redis.Options{
		Addr:     constant.Addr,     // Redis 服务器地址和端口
		Password: constant.Password, // Redis 服务器密码，如果没有设置密码则为空字符串
		DB:       constant.DB,       // Redis 数据库索引，默认为 0
	})
}
