package main

import (
	"demotest/douyin-api/globalinit/constant"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"interaction-service/model"
)

const MySQLDSN = constant.MySQL_DSN

func main() {

	// 连接数据库
	db, err := gorm.Open(mysql.Open(MySQLDSN))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	// 生成实例
	g := gen.NewGenerator(gen.Config{
		// 相对执行`go run`时的路径, 会自动创建目录
		OutPath: "./dao",

		Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	// 设置目标 db
	g.UseDB(db)

	// 创建模型的方法,生成文件在 dao 目录; 先创建结果不会被后创建的覆盖
	g.ApplyBasic(model.Comment{}, model.Favorite{})

	g.Execute()
}
