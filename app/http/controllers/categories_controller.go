package controllers

import (
	"fmt"
	"net/http"

	"github.com/StubbornYouth/goblog/app/models/article"
	"github.com/StubbornYouth/goblog/app/models/category"
	"github.com/StubbornYouth/goblog/app/requests"
	"github.com/StubbornYouth/goblog/pkg/flash"
	"github.com/StubbornYouth/goblog/pkg/route"
	"github.com/StubbornYouth/goblog/pkg/view"
)

type CategoriesController struct {
	BaseController
}

// 文章分类创建页面
func (*CategoriesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

// 文章创建提交页面
func (*CategoriesController) Store(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	_category := category.Category{
		Name: name,
	}

	errors := requests.ValidateCatrgoryForm(_category)

	if len(errors) == 0 {
		_category.Create()
		if _category.ID > 0 {
			flash.Success("文章分类创建成功")
			homeURL := route.RouteNameToURL("home")
			http.Redirect(w, r, homeURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章分类失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Category": _category,
			"Errors":   errors,
		}, "categories.create")
	}
}

// 查看分类下文章
func (cc *CategoriesController) Show(w http.ResponseWriter, r *http.Request) {
	// 获取参数id
	id := route.GetRouteVarible("id", r)

	// 读取对应分类
	_category, err := category.Get(id)

	// 获取文章结果集
	_articles, pagerData, err := article.GetByCategoryID(_category.GetStringID(), r, 2)

	if err != nil {
		cc.ResponseForSQLError(w, err, "404 文章未找到")
	} else {
		view.Render(w, view.D{"Articles": _articles, "PagerData": pagerData}, "articles.index", "articles._article_meta")
	}
}
