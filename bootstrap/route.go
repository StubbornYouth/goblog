package bootstrap

import (
	"github.com/StubbornYouth/goblog/pkg/route"
	"github.com/StubbornYouth/goblog/routes"
	"github.com/gorilla/mux"
)

// 路由初始化方法
func SetUpRoute() *mux.Router {
	Router := mux.NewRouter()
	route.SetRoute(Router)
	routes.RegisterWebRoutes(Router)
	return Router
}
