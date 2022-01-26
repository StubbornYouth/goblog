package middlewares

import (
	"net/http"

	"github.com/StubbornYouth/goblog/pkg/auth"
	"github.com/StubbornYouth/goblog/pkg/flash"
)

// 非登录用户才能访问
func Guest(next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if auth.Check() {
			flash.Warning("只有未登录用户才能访问")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next(w, r)
	}
}
