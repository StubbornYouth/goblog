package middlewares

import (
	"net/http"

	"github.com/StubbornYouth/goblog/pkg/auth"
	"github.com/StubbornYouth/goblog/pkg/flash"
)

// 登录用户才能访问
// 唯一差异的地方，是中间件书写的方式与之前不同。这是因为之前是全局中间件，而这一次我们写的中间件，可以专属指定于单个路由上
func Auth(next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.Check() {
			flash.Warning("登录用户才能访问此页面")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// 继续访问
		next(w, r)
	}
}
