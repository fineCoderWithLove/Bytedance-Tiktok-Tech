package main

import (
	"douyin/douyin-api/util"
	"fmt"
	"time"
)

func main() {
	// 连接 Redis 客户端

	//初始化操作：查询数据库并将数据存入 Redis
	util.InitRedisData()
	// 定时任务：每隔5秒将 Redis 数据写入 MySQL 数据库
	fmt.Println("redis初始化成功，定时任务开始执行")
	ticker := time.NewTicker(50 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Println("Writing Redis data to MySQL...")
		util.WriteRedisVideoToMySQL()
		util.WriteRedisUserToMySQL()
		fmt.Println("Done.")
	}
}
