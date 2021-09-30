package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestHomePage(t *testing.T) {
// 	baseURL := "http://localhost:3000"

// 	// 1. 请求 —— 模拟用户访问浏览器
// 	var (
// 		resp *http.Response
// 		err  error
// 	)
// 	resp, err = http.Get(baseURL + "/")

// 	// 2. 检测 —— 是否无错误且 200
// 	// 使用 assert.NoError() 来断言没有错误发生。第一个 t 为 testing 标准库里的 testing.T 对象，第二个参数为错误对象 err ，第三个参数为出错时显示的信息（选填）
// 	assert.NoError(t, err, "有错误发生，err 不为空")
// 	// assert.Equal() 会断言两个值相等，第一个参数同上，第二个参数是期待的状态码，第三个参数是请求返回的状态码，第四个参数为出错时显示的信息（选填）
// 	assert.Equal(t, 200, resp.StatusCode, "应返回状态码 200")
// }

// func TestAboutPage(t *testing.T) {
// 	baseURL := "http://localhost:3000"

// 	// 1. 请求 —— 模拟用户访问浏览器
// 	var (
// 		resp *http.Response
// 		err  error
// 	)
// 	resp, err = http.Get(baseURL + "/about")

// 	// 2. 检测 —— 是否无错误且 200
// 	assert.NoError(t, err, "有错误发生，err 不为空")
// 	assert.Equal(t, 200, resp.StatusCode, "应返回状态码 200")
// }

func TestAllPages(t *testing.T) {

	baseURL := "http://localhost:3000"

	// 1. 声明加初始化测试数据
	var tests = []struct {
		method   string // 请求方法
		url      string // URI
		expected int    // 状态码
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/2", 200},
		{"GET", "/articles/2/edit", 200},
		{"POST", "/articles/2", 200}, // 模拟表单提交的话，如不提供数据，会返回 200 状态码
		{"POST", "/articles", 200},
		{"POST", "/articles/1/delete", 404}, // 删除文章如果不存在，也会返回 404，为了不污染开发环境的测试，我们只需要测试到这个链接是可以返回，且可以正常处理逻辑的即可
	}

	// 2. 遍历所有测试
	for _, test := range tests {
		t.Logf("当前请求 URL: %v \n", test.url) //标准库里的辅助方法 终端打印数据
		var (
			resp *http.Response
			err  error
		)
		// 2.1 请求以获取响应
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}
		// 2.2 断言
		assert.NoError(t, err, "请求 "+test.url+" 时报错")
		assert.Equal(t, test.expected, resp.StatusCode, test.url+" 应返回状态码 "+strconv.Itoa(test.expected))
	}
}

// go test ./tests -v -count=1 增加count参数 清楚缓存 同步项目更新内容  此参数一般用以设置测试运行的次数，如果设置为 2 的话就会运行测试两次
// 可以在vscode 设置文件中给 扩展追加count
