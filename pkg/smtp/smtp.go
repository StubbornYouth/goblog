package smtp

import (
	"github.com/StubbornYouth/goblog/pkg/logger"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, title string, body string) (err error) {
	host := "smtp.qq.com"
	port := 465
	from := "335648922@qq.com"
	pwd := "qwwwjjyovhepcaih"

	// 对象初始化
	m := gomail.NewMessage()
	// 发送方
	m.SetHeader("From", from)
	// 接收方
	m.SetHeader("To", to)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	// 邮件标题
	m.SetHeader("Subject", title)
	// 邮件内容
	m.SetBody("text/html", body)
	// 附件
	// m.Attach("/home/Alex/lolcat.jpg")
	// 连接smtp
	d := gomail.NewDialer(host, port, from, pwd)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logger.LogError(err)
	}

	return err
}
