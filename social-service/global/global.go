package global

import (
	"google.golang.org/grpc"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 定义全局变量
var (
	DB *gorm.DB //mysql
	RS *redis.Client //redis
	conn *grpc.ClientConn //grpc
)


func init() {
	/*
		mysql的全局连接
	*/
	dsn := "root:0927@tcp(43.143.80.216:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
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
		Addr:     "43.143.44.118:9898", // Redis 服务器地址和端口
		Password: "192047",             // Redis 服务器密码，如果没有设置密码则为空字符串
		DB:       0,                    // Redis 数据库索引，默认为 0
	})
}
