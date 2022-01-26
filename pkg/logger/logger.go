package logger

import "log"

// 存在错误时 记录错误日志 重构checkError
func LogError(err error) {
	if err != nil {
		// 该方法打印数据以后程序就退出了，这跟我们的预期不一致，当存在错误时，我们希望记录下来，然后程序继续执行
		// log.Fatal(err)

		// log.Println() 会在 log.Print() 的基础上增加回车换行符。ln 是 line 的简写
		log.Println(err)
	}
}
