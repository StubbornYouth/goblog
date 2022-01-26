package view

import (
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/StubbornYouth/goblog/app/models/category"
	"github.com/StubbornYouth/goblog/pkg/auth"
	"github.com/StubbornYouth/goblog/pkg/flash"
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/route"
)

// D 是 map[string]interface{} 的简写
// 注意这里我们再次使用万能类型 interface{}，这是一个很常见的使用场景 更加灵活、通用
// 使用时：
// view.D{
//     "Title":  title,
//     "Body":   body,
//     "Errors": errors,
// }Copy

// 等同于：
// ArticlesFormData{
//     Title:  title,
//     Body:   body,
//     Errors: errors,
// }
type D map[string]interface{}

// 渲染通用视图
// func Render(w io.Writer, data interface{}, tplFiles ...string) {
func Render(w io.Writer, data D, tplFiles ...string) {
	RenderTemplate(w, "app", data, tplFiles...)
}

// 渲染简单视图
// func RenderSimple(w io.Writer, data interface{}, tplFiles ...string) {
func RenderSimple(w io.Writer, data D, tplFiles ...string) {
	RenderTemplate(w, "simple", data, tplFiles...)
}

// 多个模板传参
// func Render(w io.Writer, name string, article interface{}) {
// 修改后的 Render 方法支持 tplFiles... 不限参数传参，需要多少个模板，直接作为参数附加即可
// func RenderTemplate(w io.Writer, name string, data interface{}, tplFiles ...string) {
func RenderTemplate(w io.Writer, name string, data D, tplFiles ...string) {
	// // 1 设置模板相对路径
	// viewDir := "resources/views/"

	// // 遍历传参文件列表 Slice，设置正确的路径，支持 dir.filename 语法糖
	// for i, f := range tplFiles {
	// 	// 语法糖 将路由 . 转化为 / -1代表全部替换
	// 	tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	// }
	// // name = strings.Replace(name, ".", "/", -1)

	// // 所有布局模板文件 Slice
	// files, err := filepath.Glob(viewDir + "/layouts/*gohtml")
	// logger.LogError(err)

	// // 在 Slice 里新增我们的目标文件
	// // newFiles := append(files, viewDir+name+".gohtml")
	// newFiles := append(files, tplFiles...)
	// 通用模板数据
	var err error
	data["isLogined"] = auth.Check()
	data["loginUser"] = auth.User
	data["flash"] = flash.All()
	data["Categories"], err = category.All()

	newFiles := getTemplateFiles(tplFiles...)

	// 解析所有模板文件
	// tmpl, err := template.New(name + ".gohtml").Funcs(template.FuncMap{
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"RouteNameToURL": route.RouteNameToURL,
	}).ParseFiles(newFiles...)

	logger.LogError(err)

	// 渲染模板文件
	tmpl.ExecuteTemplate(w, name, data)
}

func getTemplateFiles(tplFiles ...string) []string {
	// 1 设置模板相对路径
	viewDir := "resources/views/"

	// 遍历传参文件列表 Slice，设置正确的路径，支持 dir.filename 语法糖
	for i, f := range tplFiles {
		// 语法糖 将路由 . 转化为 / -1代表全部替换
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	// name = strings.Replace(name, ".", "/", -1)

	// 所有布局模板文件 Slice
	files, err := filepath.Glob(viewDir + "/layouts/*gohtml")
	logger.LogError(err)

	// 在 Slice 里新增我们的目标文件
	// newFiles := append(files, viewDir+name+".gohtml")
	return append(files, tplFiles...)
}
