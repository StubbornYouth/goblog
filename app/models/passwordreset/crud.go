package passwordreset

import (
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/model"
)

func (passwordreset *PasswordReset) Create() error {
	if err := model.DB.Create(&passwordreset).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// 通过token 获取重置记录
// func GetByToken(token string) (PasswordReset, error) {
// 	var passwordreset PasswordReset
// 	if err := model.DB.Where("token = ?", token).First(&passwordreset).Order("id desc").Error; err != nil {
// 		return passwordreset, err
// 	}

// 	return passwordreset, nil
// }

// 通过邮箱 获取重置记录
func GetByEmail(email string) (PasswordReset, error) {
	var passwordreset PasswordReset
	if err := model.DB.Where("email = ?", email).Order("id desc").First(&passwordreset).Error; err != nil {
		return passwordreset, err
	}

	return passwordreset, nil
}
