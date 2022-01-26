package middlewares

import (
	"net/http"
	"strings"
)

// 去除url地址 最后一个/
func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 判断是否首页 根路径 不去除/
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/") // 去除url最后/
		}

		// 将请求继续传递
		next.ServeHTTP(w, r)
	})
}
