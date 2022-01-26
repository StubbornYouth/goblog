package auth

import (
	"errors"

	"github.com/StubbornYouth/goblog/app/models/user"
	"github.com/StubbornYouth/goblog/pkg/session"
	"gorm.io/gorm"
)

// 获取uid
func _getUID() string {
	_uid := session.Get("uid")
	uid, ok := _uid.(string)

	if ok && len(uid) > 0 {
		return uid
	}

	return ""
}

// user 登录获取用户信息
func User() user.User {
	uid := _getUID()

	_user, err := user.Get(uid)

	if err == nil {
		return _user
	}

	return user.User{}
}

// 尝试登录 验证邮箱密码
func Attempt(email string, password string) error {
	// 根据email 获取用户
	_user, err := user.GetByEmail(email)

	// 判断用户是否存在
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("账号不存在或密码错误")
		} else {
			return errors.New("内部错误，请稍后尝试")
		}
	}

	// 验证密码
	if !_user.ComparePassword(password) {
		return errors.New("账号不存在或密码错误")
	}

	// 登录用户 保存会话
	session.Put("uid", _user.GetStringID())

	return nil
}

// Login 指定登录用户
func Login(_user user.User) {
	session.Put("uid", _user.GetStringID())
}

// 注销用户
func Loginout() {
	session.Forget("uid")
}

// 检查是否登录
func Check() bool {
	return len(_getUID()) > 0
}
