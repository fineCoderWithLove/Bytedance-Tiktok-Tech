package constant

const (
	ErrorMsg   = "网络异常~请稍后"
	SuccessMsg = "ok"
	MySQL_DSN  = "root:0927@tcp(43.143.80.216:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	Addr       = "43.143.44.118:9898" // Redis 服务器地址和端口
	Password   = "192047"             // Redis 服务器密码，如果没有设置密码则为空字符串
	DB         = 0                    // Redis 数据库索引，默认为 0
	WordFilterFilePath     = "/home/runner/app/pub_sms_banned_words.txt"
	RedisUpdateTime        = 5//秒数
	MaxConcurrency         = 5000
)
