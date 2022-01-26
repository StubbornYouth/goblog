package passwordreset

import (
	"github.com/StubbornYouth/goblog/pkg/password"
	"gorm.io/gorm"
)

// BeforeSave GORM 的模型钩子，在保存和更新模型前调用
func (reset *PasswordReset) BeforeSave(tx *gorm.DB) (err error) {

	if !password.IsHashed(reset.Token) {
		reset.Token = password.Hash(reset.Token)
	}
	return
}
