package controllers

import (
	"fmt"
	"net/http"

	"github.com/StubbornYouth/goblog/app/models/article"
	"github.com/StubbornYouth/goblog/app/models/user"
	"github.com/StubbornYouth/goblog/pkg/route"
	"github.com/StubbornYouth/goblog/pkg/view"
)

type UserController struct {
	BaseController
}

func (uc *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVarible("id", r)

	_user, err := user.Get(id)

	if err != nil {
		uc.ResponseForSQLError(w, err, "404 未找到当前用户")
	} else {
		_articles, err := article.GetByUserID(_user.GetStringID())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		} else {
			view.Render(w, view.D{"Articles": _articles}, "articles.index", "articles._article_meta")
		}
	}

}
