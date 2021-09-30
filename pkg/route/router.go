package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// 不在当前文件初始化路由的原因 : 因为其它文件也会引用 例如控制器文件 导致与main中引用重复 出现循环引用问题
// // Router 路由对象
// var Router *mux.Router

// // Initialize 初始化路由
// func Initialize() {
// 	Router = mux.NewRouter()
// 	routes.RegisterWebRoutes(Router)
// }
var route *mux.Router

// 初始化路由对象时调用当前方法 传入路由对象
func SetRoute(r *mux.Router) {
	route = r
}

// 重构 根据路由名称获取地址
func RouteNameToURL(routeName string, pair ...string) string {
	url, err := route.Get(routeName).URL(pair...)

	if err != nil {
		// checkError(err)
		return ""
	}

	return url.String()
}

// 重构 获取路由参数
func GetRouteVarible(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)

	return vars[parameterName]
}
