package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
	"unicode/utf8"

	"github.com/StubbornYouth/goblog/app/models/article"
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/route"
	"github.com/StubbornYouth/goblog/pkg/types"
	"gorm.io/gorm"
)

type ArticlesController struct {
}

type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVarible("id", r)

	article, err := article.Get(id)

	if err != nil {
		// 判断是否未查询到数据
		// 当 Scan() 发现没有返回数据的话，会返回 sql.ErrNoRows 类型的错误
		// 修改使用gorm的方法判断
		// if err == sql.ErrNoRows {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)

			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		viewDir := "resources/views"

		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logger.LogError(err)

		newFiles := append(files, viewDir+"/articles/show.gohtml")

		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteNameToURL": route.RouteNameToURL,
			"Int64ToString":  types.Int64ToString,
		}).ParseFiles(newFiles...)
		logger.LogError(err)

		tmpl.ExecuteTemplate(w, "app", article)
	}
}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// 加载模板 只适合渲染一个文件
		// tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		// logger.LogError(err)
		// // 渲染模板
		// tmpl.Execute(w, articles)

		// 多个模板文件渲染
		// 设置模板相对路径
		viewDir := "resources/views"
		// 所有布局文件slice
		// filepath.Glob() 这是我们第一次使用 filepath 包，此包是 Go 提供的统一不同系统的路径处理包。Glob() 方法会生成与传参匹配的文件名称 Slice
		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logger.LogError(err)

		// 在slice 增加目标文件
		newFiles := append(files, viewDir+"/articles/index.gohtml")

		// 解析模板文件
		// template.ParseFiles(newFiles...) 的 ParseFiles() 是可变参数方法，三个点是 Go 提供的语法糖。
		// Slice 后加三个点，可以自动将 Slice 分解，并作为可变函数的参数
		// 以下代码等同 tmpl, err := template.ParseFiles("g.txt", "h.txt", "i.txt")
		tmpl, err := template.ParseFiles(newFiles...)
		logger.LogError(err)
		// 第一个参数和最后一个参数与 tmpl.Execute() 方法一致。中间参数 name 是最终我们想要渲染的模板名称。
		// 注意这里是模板关键词 define 定义的模板名称，不是模板文件名称
		tmpl.ExecuteTemplate(w, "app", articles)
	}
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.RouteNameToURL("articles.store")

	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, data)
}

func validateArticleFormData(title string, body string) map[string]string {
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

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		// 创建成功后 _article 对象会 返回插入的ID
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		storeUrl := route.RouteNameToURL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeUrl,
			Errors: errors,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVarible("id", r)

	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		updateURL := route.RouteNameToURL("articles.update", "id", id)

		data := ArticlesFormData{
			Title:  _article.Title,
			Body:   _article.Body,
			URL:    updateURL,
			Errors: nil,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")

		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVarible("id", r)

	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")
		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {
			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器错误")
			} else {
				if rowsAffected > 0 { // 修改成功
					showURL := route.RouteNameToURL("articles.show", "id", id)
					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					fmt.Fprint(w, "您没有做任何修改")
				}
			}
		} else {
			updateURL := route.RouteNameToURL("articles.update", "id", id)

			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}

			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)

			tmpl.Execute(w, data)
		}

	}
}

func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVarible("id", r)

	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		rowsAffected, err := _article.Delete()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		} else {
			if rowsAffected > 0 {
				indexURL := route.RouteNameToURL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}
