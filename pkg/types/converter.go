package types

import (
	"strconv"

	"github.com/StubbornYouth/goblog/pkg/logger"
)

// 存放类型转换相关方法

// 重构int64转string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// 字符串转int
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)

	if err != nil {
		logger.LogError(err)
	}

	return i
}
