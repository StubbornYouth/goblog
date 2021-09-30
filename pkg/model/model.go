package model

import (
	"github.com/StubbornYouth/goblog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// 名称重复 需要重命名
	gormlogger "gorm.io/gorm/logger"
)

// 基础模型

// DB gorm DB 对象
var DB *gorm.DB

// ConnectDB 初始化模型
func ConnectDB() *gorm.DB {

	var err error

	config := mysql.New(mysql.Config{
		DSN: "root:root@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	})

	// 准备数据库连接池
	DB, err = gorm.Open(config, &gorm.Config{
		// LogMode 里填写的是日志级别，分别如下：
		// Silent ——  静默模式，不打印任何信息 Error —— 发生错误了才打印 Warn —— 发生警告级别以上的错误才打印 Info —— 打印所有信息，包括 SQL 语句
		// 默认使用的是 Warn ，我们将其改为 Info
		// 日常开发，日志级别为 Warn 即可，否则命令太多信息会反而容易让我们错过重要的信息
		// Logger: gormlogger.Default.LogMode(gormlogger.Info),
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})

	logger.LogError(err)

	return DB
}
