package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandlerExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
			"<p>如有疑惑，请联系我们。</p>")
	}
}

func aboutHandlerExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func mainServeMux() {
	// http.HandleFunc("/", defaultHandler) // 是对DefaultServeMux.HandleFunc()的封装
	// http.HandleFunc("/about", aboutHandler)
	// http.ListenAndServe(":3000", nil) // handler 通常为 nil，此种情况下会使用 DefaultServeMux

	// 重构
	router := http.NewServeMux()

	router.HandleFunc("/", defaultHandlerExample)
	router.HandleFunc("/about", aboutHandlerExample)

	// http.ServeMux 局限性  1.不支持 URI 路径参数 2.不支持请求方法过滤 3.不支持路由命名(复用性差，如果地址发生改变 用到该地址的所有页面都需要修改)
	// 第一点和第二点可以通过代码解决 但是增加了代码维护成本

	// 路径参数处理
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 3)[2] // 对字符串通过 关键子字符串进行分割，分割成num个子字符串
		fmt.Fprint(w, "文章 ID："+id)
	})

	// 请求方法过滤
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w, "访问文章列表")
		case "POST":
			fmt.Fprint(w, "创建新的文章")
		}
	})

	http.ListenAndServe(":3000", router)

	// httpServeMux 优缺点
	// 优点
	// 1. 标准库 随着go打包安装 无需另行安装
	// 2. 测试充分
	// 3. 稳定 兼容性强
	// 4. 简单 高效
	// 缺点
	// 1.缺少web常见特性
	// 2.复杂的项目使用 需要写更多的代码

	// 标准库并不一定是最好的，有些第三方库的效率可能更高
}
