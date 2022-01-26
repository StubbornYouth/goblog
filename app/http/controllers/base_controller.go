package controllers

import (
	"fmt"
	"net/http"

	"github.com/StubbornYouth/goblog/pkg/flash"
	"github.com/StubbornYouth/goblog/pkg/logger"
	"gorm.io/gorm"
)

type BaseController struct {
}

func (bc *BaseController) ResponseForSQLError(w http.ResponseWriter, err error, message string) {
	if err == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprint(w, message)
	} else {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器错误")
	}
}

// 处理未授权的访问
func (bc *BaseController) ResponseForUnauthorized(w http.ResponseWriter, r *http.Request) {
	flash.Warning("未授权操作！")
	http.Redirect(w, r, "/", http.StatusFound)
}
