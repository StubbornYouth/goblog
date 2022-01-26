package requests

import (
	"github.com/StubbornYouth/goblog/app/models/passwordreset"

	"github.com/thedevsaddam/govalidator"
)

// 注册表单验证 go是强类型语言
func ValidarePasswordForm(data passwordreset.PasswordReset) map[string][]string {
	// 表单规则
	rules := govalidator.MapData{
		// alpha_num 只允许英文字母和数字混合
		"email": []string{"required", "min:4", "max:30", "email"},
	}

	// 定制错误消息
	messages := govalidator.MapData{
		"email": []string{
			"required:Email为必填项",
			"min:长度必须大于4",
			"max:长度必须小于30",
			"email:格式错误，请提供有效的邮箱地址",
		},
	}

	// 配置选项
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",  // Struct 标签标识符
		Messages:      messages, // 增加自定义验证提示
	}

	// 开始验证
	err := govalidator.New(opts).ValidateStruct()

	return err
}
