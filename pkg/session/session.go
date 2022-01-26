package session

import (
	"net/http"

	"github.com/StubbornYouth/goblog/pkg/config"
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/gorilla/sessions"
)

// Store gorilla sessions 的存储库
// 调用 NewCookieStore 时传参的是一串随机字符串，作为最佳实践，这个随机字串应该放置于配置中，并且随着程序的不同环境而不一致
// var Store = sessions.NewCookieStore([]byte("33446a9dcf9ea060a0a6532b166da32f304af0de"))
var Store = sessions.NewCookieStore([]byte(config.GetString("app.key")))

// Session 当前会话
var Session *sessions.Session

// 用以获取会话
var Request *http.Request

// 用以写入会话
var Response http.ResponseWriter

// StartSession 初始化会话 在中间件调用
func StartSession(w http.ResponseWriter, r *http.Request) {
	var err error
	// Store.Get 第二个参数是 Cookie的名称
	// gorilla/sessions 支持多会话，本项目我们只使用单一会话即可
	// Session, err = Store.Get(r, "goblog-session")
	Session, err = Store.Get(r, config.GetString("session.session_name"))

	logger.LogError(err)
	Request = r
	Response = w
}

// Save 保持会话
func Save() {
	// 非 HTTPS 的链接无法使用 Secure 和 HttpOnly，浏览器会报错
	// Session.Options.Secure = true
	// Session.Options.HttpOnly = true
	err := Session.Save(Request, Response)
	logger.LogError(err)
}

// 设置会话项数据
func Put(key string, value interface{}) {
	Session.Values[key] = value
	Save()
}

// Get 获取会话数据，获取数据时请做类型检测
func Get(key string) interface{} {
	return Session.Values[key]
}

// 删除某个会话项
func Forget(key string) {
	delete(Session.Values, key)
	Save()
}

// 删除当前会话
func Flush() {
	Session.Options.MaxAge = -1
	Save()
}
