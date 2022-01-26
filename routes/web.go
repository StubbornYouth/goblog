package routes

import (
	"net/http"

	"github.com/StubbornYouth/goblog/app/http/controllers"
	"github.com/StubbornYouth/goblog/app/http/middlewares"
	"github.com/gorilla/mux"
)

// 注册web 路由
func RegisterWebRoutes(r *mux.Router) {
	pc := new(controllers.PagesController)

	// r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	ac := new(controllers.ArticlesController)
	r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles", middlewares.Auth(ac.Store)).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create", middlewares.Auth(ac.Create)).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", middlewares.Auth(ac.Edit)).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(ac.Update)).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", middlewares.Auth(ac.Delete)).Methods("POST").Name("articles.delete")

	auc := new(controllers.AuthController)
	r.HandleFunc("/auth/register", middlewares.Guest(auc.Register)).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", middlewares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")
	r.HandleFunc("/auth/login", middlewares.Guest(auc.Login)).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/dologin", middlewares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")
	// 退出操作必须使用 POST 方法。恶意用户在网站上伪造图片链接 src="<你的退出链接>，将导致用户在不知情的情况下就退出登录，使用 POST 方法即可避免
	r.HandleFunc("/auth/logout", middlewares.Auth(auc.Logout)).Methods("POST").Name("auth.logout")
	// 忘记密码
	r.HandleFunc("/auth/forget", auc.Forget).Methods("GET").Name("auth.forget")
	r.HandleFunc("/auth/doforget", auc.DoForget).Methods("POST").Name("auth.doforget")
	r.HandleFunc("/auth/reset/{token}", auc.Reset).Methods("GET").Name("auth.reset")
	r.HandleFunc("/auth/doreset", auc.DoReset).Methods("POST").Name("auth.doreset")

	uc := new(controllers.UserController)
	r.HandleFunc("/users/{id:[0-9]+}", uc.Show).Methods("GET").Name("users.show")

	// 文章分类
	cc := new(controllers.CategoriesController)
	r.HandleFunc("/categories/create", middlewares.Auth(cc.Create)).Methods("GET").Name("categories.create")
	r.HandleFunc("/categories", middlewares.Auth(cc.Store)).Methods("POST").Name("categories.store")
	r.HandleFunc("/categories/{id:[0-9]+}", middlewares.Auth(cc.Show)).Methods("GET").Name("categories.show")

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

	// 全局中间件
	// 开启会话
	r.Use(middlewares.StartSession)
}
