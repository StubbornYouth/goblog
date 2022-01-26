package flash

import (
	"encoding/gob"

	"github.com/StubbornYouth/goblog/pkg/session"
)

// Flash 消息在 Web 开发中很常见，gorilla/sessions 库也对此功能提供了支持。但是内置的格式并不符合我们的要求
// Flashes Flash 消息数组类型，用以在会话中存储 map
type Flashes map[string]interface{}

var flashKey = "_flashes"

func init() {
	// 在 gorilla/sessions 中存储 map 和 struct 数据需
	// 要提前注册 gob，方便后续 gob 序列化编码、解码
	// 标准库 gob 是 Go 专属的编解码方式，是标准库自带的一个数据结构序列化的编码 / 解码工具。
	// 类似于 JSON 或 XML，不过执行效率比他们更高。特别适合在 Go 语言程序间传递数据
	gob.Register(Flashes{})
}

// 私有方法 新增一条消息提示
func addFalsh(key string, message string) {
	flashes := Flashes{}
	flashes[key] = message
	session.Put(flashKey, flashes)
	session.Save()
}

// 添加Info 类型消息提示
func Info(message string) {
	addFalsh("info", message)
}

// 添加Waring 类型消息提示
func Warning(message string) {
	addFalsh("warning", message)
}

// 添加Danger 类型消息提示
func Danger(message string) {
	addFalsh("danger", message)
}

// 添加Success 类型消息提示
func Success(message string) {
	addFalsh("success", message)
}

// All 获取所有消息
func All() Flashes {
	val := session.Get(flashKey)
	// 读取是必须做类型检测
	flashMessages, ok := val.(Flashes)
	if !ok {
		return nil
	}
	// 读取即销毁，直接删除
	session.Forget(flashKey)
	return flashMessages
}
