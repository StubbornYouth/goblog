package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/StubbornYouth/goblog/app/http/middlewares"
	"github.com/StubbornYouth/goblog/bootstrap"
	"github.com/StubbornYouth/goblog/pkg/database"

	// _ "github.com/go-sql-driver/mysql" // 匿名导入 因为这里只是引入数据库引擎 并不使用包里的方法 如果不匿名引用 代码报错不会进行编译
	"github.com/gorilla/mux"
)

// 声明包级别变量 使用var 与函数级别变量区分
// var router = mux.NewRouter()
var router *mux.Router

// 设置包级别数据库 *sql.DB 结构体实例 方便使用
var db *sql.DB

func main() {
	database.Initialize()
	db = database.DB

	// route.Initialize()
	// router = route.Router
	// 初始化 gorm DB
	bootstrap.SetupDB()
	router = bootstrap.SetUpRoute()

	// 通过路由获取url
	homeUrl, _ := router.Get("home").URL()
	fmt.Println("homeUrl：", homeUrl)

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router)) // removeTrailingSlash对路由进行处理
}
