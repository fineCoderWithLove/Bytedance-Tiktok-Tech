package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// Conf 全局配置变量
var Conf = new(config)

type config struct {
	Mysql  *MysqlConfig `mapstructure:"mysql" json:"mysql"`
	System *System      `mapstructure:"system" json:"system"`
}

func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}

	// 设置文件路径
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "./")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s \n", err))
	}

	// 将viper 中的json 信息映射到 全局Conf中
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}

}

// MysqlConfig 对应yml的中 mysql
type MysqlConfig struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}

// System 对应yml中的system
type System struct {
	Mode          string `mapstructure:"mode" json:"mode"`
	UrlPathPrefix string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Host          string `mapstructure:"host" json:"host"`
	Port          int    `mapstructure:"port" json:"port"`
}
