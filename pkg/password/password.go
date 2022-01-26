package password

import (
	"github.com/StubbornYouth/goblog/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// Hash 使用bcrypt对密码进行加密
func Hash(password string) string {
	// 第二个参数是cost值 越大耗费时间越长 建议大于12
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.LogError(err)

	return string(bytes)
}

// CheckHash 对比明文密码和数据库的哈希值
func CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	logger.LogError(err)
	return err == nil
}

// 判断字符串是否已加密
func IsHashed(str string) bool {
	// bcrypt加密国的字符串长度为60
	return len(str) == 60
}
