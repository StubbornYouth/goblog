package bootstrap

import (
	"time"

	"github.com/StubbornYouth/goblog/app/models/category"
	"github.com/StubbornYouth/goblog/app/models/passwordreset"
	"github.com/StubbornYouth/goblog/app/models/user"

	"github.com/StubbornYouth/goblog/app/models/article"
	"github.com/StubbornYouth/goblog/pkg/config"
	"github.com/StubbornYouth/goblog/pkg/model"
	"gorm.io/gorm"
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
	// sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxOpenConns(config.GetInt("database.mysql.max_idle_connections"))
	// 设置最大空闲连接数
	// sqlDB.SetMaxIdleConns(25)
	sqlDB.SetMaxIdleConns(config.GetInt("database.mysql.max_open_connections"))

	// 设置连接断开时间
	// sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second)

	// 创建维护数据表结构
	migration(db)
}

// 数据库迁移
// GORM 自带了自动迁移功能，可以针对设置的模型 Struct 来自动创建数据表结构。免去了我们手动维护 SQL 的烦恼，自动迁移也有统一多个数据库系统的好处。
// 使用自动迁移很简单，只需要调用 AutoMigrate() 方法并将数据模型 Struct 传参进去即可。
// 因为这是一个全局动作，我们这个操作放置于数据库初始化的地方
func migration(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
		&passwordreset.PasswordReset{},
		&category.Category{},
	)
}
