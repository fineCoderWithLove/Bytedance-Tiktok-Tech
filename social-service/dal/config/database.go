package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		Conf.Mysql.Username,
		Conf.Mysql.Password,
		Conf.Mysql.Host,
		Conf.Mysql.Port,
		Conf.Mysql.Database,
		Conf.Mysql.Charset,
		Conf.Mysql.Collation,
		Conf.Mysql.Query,
	)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Printf("数据库连接错误：%s", err)
		return
	}

	// 全局赋值
	DB = db

	log.Printf("数据库连接成功")
}
