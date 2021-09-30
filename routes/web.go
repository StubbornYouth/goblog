package routes

import (
	"net/http"

	"github.com/StubbornYouth/goblog/app/http/controllers"
	"github.com/gorilla/mux"
)

// 注册web 路由
func RegisterWebRoutes(r *mux.Router) {
	pc := new(controllers.PagesController)

	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	ac := new(controllers.ArticlesController)
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	// 404
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	// 静态资源
	// PathPrefix() 匹配参数里 /css/ 前缀的 URI
	// 链式调用 Handler() 指定处理器为 http.FileServer()
	// http.FileServer() 是文件目录处理器，参数 http.Dir("./public") 是指定在此目录下寻找文件
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// 中间件：强制内容类型为 HTML
	// Execute() 在执行时会设置正确的 HTML 标头 解析静态文件所用到的 http.FileServer() 内部也会根据文件后缀设置正确的标头 所以标头这块不需要我们干预
	// r.Use(middlewares.ForceHTML)
}
