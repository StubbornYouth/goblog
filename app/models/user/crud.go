package user

import (
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/model"
	"github.com/StubbornYouth/goblog/pkg/password"
)

func (user *User) Create() error {
	if err := model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// 通过id 获取用户信息
func Get(idstr string) (User, error) {
	var user User
	if err := model.DB.First(&user, idstr).Error; err != nil {
		return user, err
	}

	return user, nil
}

// 通过email 获取用户信息
func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// 更新密码
func UpdatePassword(_user User) (RowsAffected int64, err error) {
	if !password.IsHashed(_user.Password) {
		_user.Password = password.Hash(_user.Password)
	}
	result := model.DB.Model(&User{}).Where("email = ?", _user.Email).Update("password", _user.Password)

	// result 返回两个参数 result.rowsAffected更新变化条数 result.Error 错误信息
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil

}
