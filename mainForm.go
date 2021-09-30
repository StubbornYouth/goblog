package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// 声明包级别变量 使用var 与函数级别变量区分
var routerForm = mux.NewRouter()

func homeFormHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html;charset=utf8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
}

func aboutFormHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
		"<p>如有疑惑，请联系我们。</p>")
}

func articlesShowFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 获取路由参数
	id := vars["id"]
	fmt.Fprint(w, "文章ID："+id)
}

func articlesIndexFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "查看文章列表")
}

// 构建struct 用以给模板文件传输变量使用
type ArticlesFormDataForm struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesStoreFormHandler(w http.ResponseWriter, r *http.Request) {
	/** 表单数据获取 */
	// 从请求中解析参数 必须执行这段代码 否则 r.PostForm 和 r.Form 获取到的为空数组
	// error := r.ParseForm()

	// 判断解析是否正确
	// if error != nil {
	// 	fmt.Fprintf(w, "请提交正确的数据")
	// 	return
	// }

	// r.Form 存储了post put get 参数 r.PostForm 存储了 post put参数
	// title := r.PostForm.Get("title")

	// fmt.Fprintf(w, "POST PostForm %v <br>", r.PostForm)
	// fmt.Fprintf(w, "POST Form %v <br>", r.Form)
	// fmt.Fprintf(w, "title 的值为 %v", title)

	// 不想获取请求的所有内容 可以逐个获取 不需要使用 r.ParseForm
	// fmt.Fprintf(w, "POST Form中标题为 %v <br>", r.FormValue("title"))
	// fmt.Fprintf(w, "POST PostForm中标题为 %v <br>", r.PostFormValue("title"))
	// fmt.Fprintf(w, "POST Form中内容为 %v <br>", r.FormValue("body"))
	// fmt.Fprintf(w, "POST PostForm中内容为 %v <br>", r.PostFormValue("body"))

	/** 表单验证 */
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	// 声明一个map集合变量 string 数据类型
	errors := make(map[string]string)

	// len 获取字节数 一个utf-8中文字符占3个字节
	// 需要获取实际字符长度 使用utf8包中的utf8.RuneCountInString()来计数
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度必须介于3~40个字符"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容不能少于10个字符"
	}

	// 判断是否有错误发生
	if len(errors) == 0 {
		fmt.Fprintf(w, "title值为：%v<br>", title)
		fmt.Fprintf(w, "title字节长度为：%v<br>", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body值为：%v<br>", body)
		fmt.Fprintf(w, "body字节长度为：%v<br>", utf8.RuneCountInString(body))
	} else {
		// fmt.Fprintf(w, "有错误发生 errors的值为：%v<br>", errors)
		// 处理错误提示
		// html := `
		// <!DOCTYPE html>
		// <html lang="en">
		// <head>
		// 	<title>创建文章 —— 我的技术博客</title>
		// 	<style type="text/css">.error {color: red;}</style>
		// </head>
		// <body>
		// 	<form action="{{ .URL }}" method="post">
		// 		<p><input type="text" name="title" value="{{ .Title }}"></p>
		// 		{{ with .Errors.title }}
		// 		<p class="error">{{ . }}</p>
		// 		{{ end }}
		// 		<p><textarea name="body" cols="30" rows="10">{{ .Body }}</textarea></p>
		// 		{{ with .Errors.body }}
		// 		<p class="error">{{ . }}</p>
		// 		{{ end }}
		// 		<p><button type="submit">提交</button></p>
		// 	</form>
		// </body>
		// </html>
		// `

		storeURL, _ := routerForm.Get("articles.store").URL()

		// 构建 ArticlesFormData 数据
		data := ArticlesFormDataForm{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}

		// template.New() 包的初始化
		// tmpl, err := template.New("create-form").Parse(html)
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

// 设置表头中间件
func forceHTMLFormMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置表头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 去除url地址 最后一个/
func removeTrailingFormSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 判断是否首页 根路径 不去除/
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/") // 去除url最后/
		}

		// 将请求继续传递
		next.ServeHTTP(w, r)
	})
}

// 文章新增表单
func articlesCreateFormHandle(w http.ResponseWriter, r *http.Request) {
	// 上标符号 ` 来书写 HTML 代码，一般多行字符串可以使用这种方式
	// html := `
	// <!DOCTYPE html>
	// <html lang="en">
	// <head>
	// 	<title>创建文章 —— 我的技术博客</title>
	// </head>
	// <body>
	// 	<form action="%s" method="post">
	// 		<p><input type="text" name="title"></p>
	// 		<p><textarea name="body" cols="30" rows="10"></textarea></p>
	// 		<p><button type="submit">提交</button></p>
	// 	</form>
	// </body>
	// </html>
	// `

	storeUrl, _ := routerForm.Get("articles.store").URL()

	data := ArticlesFormDataForm{
		Title:  "",
		Body:   "",
		URL:    storeUrl,
		Errors: nil,
	}

	// 由于Go模板引擎使用{{}}作为标识 可能会和vue.js等前端框架冲突，可以修改默认标识符
	// template.New("crete-form").Delims("{[", "]}").ParseFiles("resources/views/articles/create.gohtml")
	//获取模板文件对象
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

	if err != nil {
		// 在 Go 中，一般 err 处理方式可以是给用户提示或记录到错误日志里，这种很多时候为 业务逻辑错误。当有重大错误，或者系统错误时，例如无法加载模板文件，就使用 panic()
		panic(err)
	}

	// fmt.Fprintf(w, html, storeUrl)
	tmpl.Execute(w, data)
}

func mainForm() {
	// gorilla/mux 因实现了 net/http 包的 http.Handler 接口 兼容 http.serveMux
	// gorilla/mux 精准匹配原则 路由只会匹配准确指定的规则，这个比较好理解，也是较常见的匹配方式
	// http.serveMux 长度优先匹配 一般用在静态路由上（不支持动态元素如正则和 URL 路径参数），优先匹配字符数较多的规则
	// router := mux.NewRouter()

	// 设置地址 请求方式 路由名
	routerForm.HandleFunc("/", homeFormHandler).Methods("GET").Name("home")
	routerForm.HandleFunc("/about", aboutFormHandler).Methods("GET").Name("about")
	// 正则验证id
	routerForm.HandleFunc("/articles/{id:[0-9]+}", articlesShowFormHandler).Methods("GET").Name("articles.show")
	routerForm.HandleFunc("/articles", articlesIndexFormHandler).Methods("GET").Name("articles.index")
	routerForm.HandleFunc("/articles", articlesStoreFormHandler).Methods("POST").Name("articles.store")
	routerForm.HandleFunc("/articles/create", articlesCreateFormHandle).Methods("GET").Name("articles.create")

	// 自定义404
	routerForm.NotFoundHandler = http.HandlerFunc(notFoundFormHandler)

	// 强制路由header中间件
	routerForm.Use(forceHTMLFormMiddleware)

	// 通过路由获取url
	homeUrl, _ := routerForm.Get("home").URL()
	fmt.Println("homeUrl：", homeUrl)
	articlesUrl, _ := routerForm.Get("articles.show").URL("id", "23")
	fmt.Println("articlesUrl：", articlesUrl)

	http.ListenAndServe(":3000", removeTrailingFormSlash(routerForm)) // removeTrailingSlash对路由进行处理

	// 终端测试访问 在 Gorilla Mux 中，如未指定请求方法，默认会匹配所有方法
	// curl -X POST http://localhost:3000/articles post新增文章
	// curl http://localhost:3000/articles // get查看文章列表

	// url 访问中 末尾加/ 由于精准匹配原则 会404响应 无法找到页面
	// 对于这个问题 Gorilla Mux 提供了一个 StrictSlash(value bool) 函数  会将地址301重定向到正确的地址，但是301重定向是GET请求 不适用于post请求
	// router := mux.NewRouter()Copy 修改为 router := mux.NewRouter().StrictSlash(true)
	// 不能使用中间件 因为执行顺序的问题 路由匹配优先于中间件执行,使用函数进行解析处理
}
