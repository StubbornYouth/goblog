package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	// _ "github.com/go-sql-driver/mysql" // 匿名导入 因为这里只是引入数据库引擎 并不使用包里的方法 如果不匿名引用 代码报错不会进行编译
	"github.com/gorilla/mux"
)

// 声明包级别变量 使用var 与函数级别变量区分
var routerDb = mux.NewRouter()

// 设置包级别数据库 *sql.DB 结构体实例 方便使用
var dbDb *sql.DB

func initDBDb() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 准备数据库连接池 生成DSN信息
	// DSN 全称为 Data Source Name，表示 数据源信息，用于定义如何连接数据库。不同数据库的 DSN 格式是不同的，这取决于数据库驱动的实现
	// 使用 sql.Open() 函数便可初始化并返回一个 *sql.DB 结构体实例，使用 sql.Open() 函数只要传入驱动名称及对应的 DSN 便可 需要连接不同数据库时，只需修改驱动名与 DSN 即可
	db, err = sql.Open("mysql", config.FormatDSN())
	checkErrorDb(err)

	// 设置最大连接数 设连接池最大打开数据库连接数，<= 0 表示无限制，默认
	// 实验表明，在高并发的情况下，将值设为大于 10，可以获得比设置为 1 接近六倍的性能提升。而设置为 10 跟设置为 0（也就是无限制），在高并发的情况下，性能差距不明显
	// 不要超出数据库系统设置的最大连接数 mysql8 是151 可通过show variables like 'max_connections'; 数据库查询
	db.SetMaxOpenConns(100)

	// 设置最大空闲连接数 设置连接池最大空闲数据库连接数，<= 0 表示不设置空闲连接数，默认为 2
	// 实验表明，在高并发的情况下，将值设为大于 0，可以获得比设置为 0 超过 20 倍的性能提升
	db.SetMaxIdleConns(25)

	// 设置每个链接过期时间
	// 理论上来讲，在并发的情况下，此值越小，连接就会越快被关闭，也意味着更多的连接会被创建。 比较保守的做法是设置五分钟
	// SetConnMaxLifetime 要求传参的是一个 time.Duration 对象，所以这里使用了 time.Minute，这也是我们初次使用标准库里的关于处理时间的包 —— time
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接数据库
	err = db.Ping()
	checkErrorDb(err)
}

// Go 语言中根据首字母的大小写来确定可以访问的权限。无论是函数名、方法名、常量、变量名还是结构体的名称
// 如果首字母大写，则可以被其他的包访问；如果首字母小写，则只能在本包中使用。可以简单的理解成，首字母大写是公有的，首字母小写是私有的
func creteTablesDb() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
	id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	body longtext COLLATE utf8mb4_unicode_ci
	);`

	// 一般没有返回集的sql 语句使用 Exec执行 例如 insert update delete 等语句
	// func (db *DB) Exec(query string, args ...interface{}) (Result, error) Exec() 方法的第一个返回值为一个实现了 sql.Result 接口的类型，sql.Result 的定义如下
	// type Result interface {
	// 	LastInsertId() (int64, error)    // 使用 INSERT 向数据插入记录，数据表有自增 id 时，该函数有返回值
	// 	RowsAffected() (int64, error)    // 表示影响的数据表行数
	// }
	// 可以用 sql.Result 中的 LastInsertId() 方法或 RowsAffected() 来判断 SQL 语句是否执行成功。
	// 因为我们执行的是 CREATE TABLE IF NOT EXISTS 语句，会被重复执行，所以这里判断返回结果意义不大，主要判断返回的第二个参数 err 是否有问题
	_, err := db.Exec(createArticlesSQL)

	checkErrorDb(err)
}

func checkErrorDb(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func homeDbHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html;charset=utf8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
}

func aboutDbHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

// (int64, error) 函数返回值类型
func saveArticleDbToDB(title string, body string) (int64, error) {
	// 变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 获取一个prepare预处理语句 永远不要相信用户提交的数据 防止sql注入
	// stmt statement缩写 预处理返回的一个 *sql.Stmt 指针对象 会占用sql连接
	// = 是赋值 :=是声明变量并赋值
	stmt, err = db.Prepare("INSERT INTO articles (title,body) values(?,?)")

	// 错误检测
	if err != nil {
		return 0, err
	}

	// 执行sql请求，传参
	// rs sql.Result type Result interface {
	// 	LastInsertId() (int64, error)    // 使用 INSERT 向数据插入记录，数据表有自增 id 时，该函数有返回值
	// 	RowsAffected() (int64, error)    // 表示影响的数据表行数
	// }
	rs, err = stmt.Exec(title, body)

	// 错误检测
	if err != nil {
		return 0, err
	}

	// 关闭SQL连接 释放资源
	// defer 延迟执行语句，Go 语言的 defer 语句会将其后面跟随的语句进行延迟处理，在 defer 归属的函数即将返回时，执行被延迟的语句
	defer stmt.Close()

	// 根据返回的rs 自增id判断是否插入成功
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

func notFoundDbHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
		"<p>如有疑惑，请联系我们。</p>")
}

// 设置文章struct 对应一条文章数据
type ArticleDb struct {
	Title, Body string
	ID          int64
}

// Link方法生成文章链接
// 区别于不同 的函数方法 以下是声明的对比
// type Object struct {
//     ...
// }
// Object 的方法
// func (obj *Object) method() {
//     ...
// }

// 只是一个函数
// func function() {
//     ...
// }

// 调用方法：
// o := new(Object)
// o.method()

// // 调用函数
// function()
func (article ArticleDb) LinkDb() string {
	showUrl, err := router.Get("articles.show").URL("id", strconv.FormatInt(article.ID, 10))
	if err != nil {
		checkErrorDb(err)
		return ""
	}

	return showUrl.String()
}

// RouteName2URL 通过路由名称来获取 URL
func RouteNameToURLDb(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		checkErrorDb(err)
		return ""
	}

	return url.String()
}

// Int64ToString 将 int64 转换为 string
func Int64ToStringDb(num int64) string {
	return strconv.FormatInt(num, 10)
}

// 删除文章方法封装
func (article ArticleDb) DeleteDb() (rowsAffected int64, err error) {
	// query := "DELETE * FROM articles WHERE id = ?"
	// row, err := db.Exec(query, strconv.FormatInt(article.ID, 10))
	// 这里我们用的是 Exec()，一般在 CREATE/UPDATE/DELETE 时使用。这里使用的是纯文本模式的查询模式，
	// 因为 ID 我们是从数据库里拿出来的，是自增 ID ，无需担心 SQL 注入，这样可以少发送一次 SQL 请求
	row, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(article.ID, 10))
	if err != nil {
		return 0, err
	}

	if n, _ := row.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func articlesDbShowHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVaribleDb("id", r)

	// article := Article{}
	// query := "SELECT * FROM articles WHERE id = ?"

	// QueryRow 是可变参数方法 func (db *DB) QueryRow(query string, args ...interface{}) *Row
	// 参数 只有一个时 为纯文本模式 有多个参数时 为prepare模式
	// 关于 Prepare 模式和纯文本模式，这里还需提两点：
	// 1.使用 Prepare 模式会发送两个 SQL  请求到 MySQL 服务器上，而纯文本模式只有一个；
	// 2.在使用路由参数过滤只允许数字的情况下，可以放心使用纯文本模式无需担心 SQL 注入，这里有意使用 Prepare 模式是为了课程的需要
	// Scan方法为链式调用  QueryRow 方法返回的时一个指针变量，占用一个Sql连接 调用Scan方法时，会释放sql来连接，所以每次 QueryRow 后使用 Scan 是必须的
	// err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body) //多个参数下 封装了prepare方法的调用
	// 以上语句等于
	// stmt, err := db.Prepare(query)
	// checkError(err)
	// defer stmt.Close()
	// err = stmt.QueryRow(id).Scan(&article.ID, &article.Title, &article.Body)

	article, err := getArticleByIdDb(id)

	if err != nil {
		// 判断是否未查询到数据
		// 当 Scan() 发现没有返回数据的话，会返回 sql.ErrNoRows 类型的错误
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)

			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 数据库错误
			checkErrorDb(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		// fmt.Fprint(w, "获取成功，文章标题为："+article.Title)
		// 加载渲染模板
		// tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")

		// 这一次是使用 template.New() 先初始化，然后使用 Funcs() 注册函数，再使用 ParseFiles ()，
		// 需要注意的是 New() 的参数是模板名称，需要对应 ParseFiles() 中的文件名，否则会无法正确读取到模板，最终显示空白页面。
		// Funcs() 方法的传参是 template.FuncMap 类型的 Map 对象。键为模板里调用的函数名称，值为当前上下文的函数名称
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteNameToURL": RouteNameToURLDb,
			"Int64ToString":  Int64ToStringDb,
		}).ParseFiles("resources/views/articles/show.gohtml")

		checkErrorDb(err)

		tmpl.Execute(w, article)
	}
}

func articlesDbIndexHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM articles"

	// Query方法和 QueryRow 一样 有文本模式 和 prepare 模式 由于这里我们不需要获取相关参数进行查询 所以采用文本模式
	// Query() 返回结果集 Rows ，包含从数据库里读取出来的数据和 SQL 连接
	// 使用Query() 方法要注意以下几点
	// 1.在每一次 for rows.Next() 后，都记得要检测下是否有错误发生，调用 rows.Err() 可获取到错误；
	// 2.使用 rows.Next() 遍历数据，遍历到最后内部遇到 EOF 错误，会自动调用 rows.Close() 将 SQL 连接关闭；
	// 3.使用 rows.Next() 遍历时，如遇错误，SQL 连接也会自动关闭；
	// 4.rows.Close() 可调用多次，使用 rows.Close() 可保证 SQL 连接永远是关闭的。defer rows.Close() 需在检测 err 以后调用，否则会让运行时 panic ；
	// 5.牢记在获取到结果集后，必须执行 defer rows.Close()。这样做能防止有时你在函数里过早 return ，或者其他操作忘记关闭资源，这是一个值得培养的良好习惯；
	// 6.如果你在循环中执行  Query() 并获取 Rows 结果集，请不要使用 defer ，而是直接调用 rows.Close()，因为 defer 不会立刻执行，而是在函数执行结束后执行
	rows, err := db.Query(query)
	checkErrorDb(err)
	defer rows.Close()

	var articles []ArticleDb
	// 循环遍历
	for rows.Next() {
		var article ArticleDb
		// 循环遍历每一篇文章 赋值给article对象
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkErrorDb(err)
		// 将article追加到数组中
		articles = append(articles, article)
	}
	// 检测遍历时是否发生错误
	err = rows.Err()
	checkErrorDb(err)

	// 加载模板
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkErrorDb(err)
	// 渲染模板
	tmpl.Execute(w, articles)
}

// 构建struct 用以给模板文件传输变量使用
type ArticlesFormDataDb struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesDbStoreHandler(w http.ResponseWriter, r *http.Request) {
	/** 表单验证 */
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	// 表单验证函数
	errors := validateArticleFormDataDb(title, body)

	// 判断是否有错误发生
	if len(errors) == 0 {
		lastInsertID, err := saveArticleDbToDB(title, body)
		if lastInsertID > 0 {
			// FormatInt() 方法来将类型为 int64 的 lastInsertID 转换为字符串。此方法的第二个参数 10 为十进制
			fmt.Fprintf(w, "插入成功,插入数据ID："+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkErrorDb(err)
			// 设置服务器响应
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "服务器500错误")
		}
	} else {

		storeURL, _ := router.Get("articles.store").URL()

		// 构建 ArticlesFormData 数据
		data := ArticlesFormDataDb{
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
func forceDbHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置表头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 去除url地址 最后一个/
func removeTrailingSlashDb(next http.Handler) http.Handler {
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
func articlesDbCreateHandle(w http.ResponseWriter, r *http.Request) {
	storeUrl, _ := router.Get("articles.store").URL()

	data := ArticlesFormDataDb{
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

func articlesDbEditHandler(w http.ResponseWriter, r *http.Request) {
	// 获取参数
	// vars := mux.Vars(r)
	// id := vars["id"]
	id := getRouteVaribleDb("id", r)

	// 查询当前文章信息
	// article := Article{}
	// query := "SELECT * FROM articles WHERE id = ?"
	// err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	article, err := getArticleByIdDb(id)

	// 判断查询是否成功
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 数据库错误
			checkErrorDb(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		// 查询成功 引入渲染模板
		updateUrl, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormDataDb{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateUrl,
			Errors: nil,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		checkErrorDb(err)

		tmpl.Execute(w, data)
	}
}

// 代码重构 少写代码 提高可维护性
// 获取请求参数
func getRouteVaribleDb(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

// 通过id获取文章信息
func getArticleByIdDb(id string) (ArticleDb, error) {
	query := "SELECT * FROM articles WHERE id = ?"
	article := ArticleDb{}
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	return article, err
}

// 表单验证
func validateArticleFormDataDb(title string, body string) map[string]string {
	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "文章标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}
	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

func articlesDbUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVaribleDb("id", r)

	_, err := getArticleByIdDb(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkErrorDb(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		// 4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormDataDb(title, body)

		if len(errors) == 0 {
			query := "UPDATE articles SET title = ?,body = ? WHERE id = ?"
			// Exec 是DB方法 与stmt.exec区分 同样 该方法与 QueryRow 类似 支持单独参数的纯文本模式 与 多个参数的 Prepare 模式
			// 善用 Exec() 的 Prepare 模式 来防范 SQL 注入攻击 使用此方法来处理 CREATE、UPDATE、DELETE 类型的 SQL
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				checkErrorDb(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "更新失败，服务器错误")
			} else {
				// 判断数据库更新数量
				if n, _ := rs.RowsAffected(); n > 0 {
					showUrl, _ := router.Get("articles.show").URL("id", id)
					// 跳转到详情页
					http.Redirect(w, r, showUrl.String(), http.StatusFound)
				} else {
					fmt.Fprint(w, "您没有做任何修改！")
				}
			}
		} else {
			updateUrl, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormDataDb{
				Title:  title,
				Body:   body,
				URL:    updateUrl,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			checkErrorDb(err)
			tmpl.Execute(w, data)
		}
	}
}

func articlesDbDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVaribleDb("id", r)

	article, err := getArticleByIdDb(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkErrorDb(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		rowsAffected, err := article.DeleteDb()
		if err != nil {
			// 大概率数据库错误
			checkErrorDb(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		} else {
			if rowsAffected > 0 {
				// 重定向到文章首页
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

func mainDb() {
	initDBDb()
	// 创建表
	creteTablesDb()

	// gorilla/mux 因实现了 net/http 包的 http.Handler 接口 兼容 http.serveMux
	// gorilla/mux 精准匹配原则 路由只会匹配准确指定的规则，这个比较好理解，也是较常见的匹配方式
	// http.serveMux 长度优先匹配 一般用在静态路由上（不支持动态元素如正则和 URL 路径参数），优先匹配字符数较多的规则
	// router := mux.NewRouter()

	// 设置地址 请求方式 路由名
	router.HandleFunc("/", homeDbHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutDbHandler).Methods("GET").Name("about")
	// 正则验证id
	router.HandleFunc("/articles/{id:[0-9]+}", articlesDbShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesDbIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesDbStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesDbCreateHandle).Methods("GET").Name("articles.create")
	// 增加编辑更新路由
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesDbEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesDbUpdateHandler).Methods("POST").Name("articles.update")
	// 删除路由
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDbDeleteHandler).Methods("POST").Name("articles.delete")

	// 自定义404
	router.NotFoundHandler = http.HandlerFunc(notFoundDbHandler)

	// 强制路由header中间件
	router.Use(forceDbHTMLMiddleware)

	// 通过路由获取url
	homeUrl, _ := router.Get("home").URL()
	fmt.Println("homeUrl：", homeUrl)
	articlesUrl, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articlesUrl：", articlesUrl)

	http.ListenAndServe(":3000", removeTrailingSlashDb(router)) // removeTrailingSlash对路由进行处理

	// 终端测试访问 在 Gorilla Mux 中，如未指定请求方法，默认会匹配所有方法
	// curl -X POST http://localhost:3000/articles post新增文章
	// curl http://localhost:3000/articles // get查看文章列表

	// url 访问中 末尾加/ 由于精准匹配原则 会404响应 无法找到页面
	// 对于这个问题 Gorilla Mux 提供了一个 StrictSlash(value bool) 函数  会将地址301重定向到正确的地址，但是301重定向是GET请求 不适用于post请求
	// router := mux.NewRouter()Copy 修改为 router := mux.NewRouter().StrictSlash(true)
	// 不能使用中间件 因为执行顺序的问题 路由匹配优先于中间件执行,使用函数进行解析处理
}
