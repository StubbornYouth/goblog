package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func homeMuxHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html;charset=utf8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
}

func aboutMuxHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
		"<p>如有疑惑，请联系我们。</p>")
}

func articlesShowMuxHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 获取路由参数
	id := vars["id"]
	fmt.Fprint(w, "文章ID："+id)
}

func articlesIndexMuxHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "查看文章列表")
}

func articlesStoreMuxHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "新增文章")
}

// 设置表头中间件
func forceHTMLMuxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置表头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 去除url地址 最后一个/
func removeTrailingMuxSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 判断是否首页 根路径 不去除/
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/") // 去除url最后/
		}

		// 将请求继续传递
		next.ServeHTTP(w, r)
	})
}

func mainMux() {
	// gorilla/mux 因实现了 net/http 包的 http.Handler 接口 兼容 http.serveMux
	// gorilla/mux 精准匹配原则 路由只会匹配准确指定的规则，这个比较好理解，也是较常见的匹配方式
	// http.serveMux 长度优先匹配 一般用在静态路由上（不支持动态元素如正则和 URL 路径参数），优先匹配字符数较多的规则
	router := mux.NewRouter()

	// 设置地址 请求方式 路由名
	router.HandleFunc("/", homeMuxHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutMuxHandler).Methods("GET").Name("about")
	// 正则验证id
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowMuxHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexMuxHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreMuxHandler).Methods("POST").Name("articles.store")

	// 自定义404
	router.NotFoundHandler = http.HandlerFunc(notFoundMuxHandler)

	// 强制路由header中间件
	router.Use(forceHTMLMuxMiddleware)

	// 通过路由获取url
	homeUrl, _ := router.Get("home").URL()
	fmt.Println("homeUrl：", homeUrl)
	articlesUrl, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articlesUrl：", articlesUrl)

	http.ListenAndServe(":3000", removeTrailingMuxSlash(router)) // removeTrailingSlash对路由进行处理

	// 终端测试访问 在 Gorilla Mux 中，如未指定请求方法，默认会匹配所有方法
	// curl -X POST http://localhost:3000/articles post新增文章
	// curl http://localhost:3000/articles // get查看文章列表

	// url 访问中 末尾加/ 由于精准匹配原则 会404响应 无法找到页面
	// 对于这个问题 Gorilla Mux 提供了一个 StrictSlash(value bool) 函数  会将地址301重定向到正确的地址，但是301重定向是GET请求 不适用于post请求
	// router := mux.NewRouter()Copy 修改为 router := mux.NewRouter().StrictSlash(true)
	// 不能使用中间件 因为执行顺序的问题 路由匹配优先于中间件执行,使用函数进行解析处理
}
