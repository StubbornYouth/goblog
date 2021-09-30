// 每一段 Go 程序都 必须 属于一个包。而 main 包在 Go 程序中有特殊的位置。
// 如果一段程序是属于 main 包的，那么当执行 go install 或者 go run 时就会将其生成二进制文件，当执行这个文件时，就会调用 main 函数。
// main 包里的 main 函数相当于应用程序的入口。要想生成可执行的二进制文件，必须把代码写在 main 包里，而且其中必须包含一个 main 函数
package main

// Go 语言标准库是由 Go 官方团队维护，包含在 Go 语言安装包中的 Go 包  以下引入的两个包就是标准库中的包
import ( // 引入包
	"fmt"      // 使用频率很高 它是 format 的缩写，fmt 包含有格式化 I/O 函数。主要分为向外输出内容和获取输入内容两大部分
	"net/http" // 提供了 HTTP 编程有关的接口 内部封装了 TCP 连接和报文解析的复杂琐碎细节。http 提供了 HTTP 客户端和服务器实现
)

// w 为http.ResponseWriter对象简写(响应) r为http.Request对象简写(请求)
/**
* 常见操作
* 获取用户参数 r.URL.Query() r.Header.Get("User-Agent") 返回 500 状态码 w.WriteHeader(http.StatusInternalServerError) 设置返回标头 w.Header().Set("name", "my name is smallsoup")
 */
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
	// fmt.Fprint(w, "请求路径为："+r.URL.Path)
	// 设置html文档表头 html标签等能够正常渲染
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 利用/指代任意路劲 可以根据不同路径响应不同内容
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
			"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
	} else {
		// 设置响应404
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
			"<p>如有疑惑，请联系我们。</p>")
	}
}
func mainInit() {
	// http.HandleFunc 用以指定处理 HTTP 请求的函数，此函数允许我们只写一个 handler（在此例子中 handlerFunc，可任意命名），请求会通过参数传递进来，使用者只需与 http.Request 和 http.ResponseWriter 两个对象交互即可
	// 这里的/指代的是任意路径 而不是根路径
	http.HandleFunc("/", handlerFunc)
	// nil是go语言中预先的标识符 我们可以直接使用nil，而不用声明它。 nil可以代表很多类型的零值 在go语言中，nil可以代表下面这些类型的零值
	// 一个已声明的变量 但是没有赋值的话 即是nil
	http.ListenAndServe(":3000", nil) // 用以监听本地 3000 端口以提供服务，标准的 HTTP 端口是 80 端口 另一个 Web 常用是 HTTPS 的 443 端口
}

// tmp目录是 air 命令的编译文件存放地 git提交时需要忽略这部分内容 新建.gitignore文件设置
