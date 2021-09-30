package bootstrap

import (
	"time"

	"github.com/StubbornYouth/goblog/pkg/model"
)

// 初始化数据库连接和ORM
func SetupDB() {
	// 建立数据库连接池
	db := model.ConnectDB()

	// 命令行打印数据库请求的信息
	// *gorm.DB 对象有一个方法 DB() 可以直接获取到 database/sql 包里的 *sql.DB 对象。
	// 从以下代码不难看出，GORM 底层也是使用 database/sql 来管理连接池
	sqlDB, _ := db.DB()

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(25)

	// 设置连接断开时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
