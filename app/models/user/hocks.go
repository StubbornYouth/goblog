package user

import (
	"github.com/StubbornYouth/goblog/pkg/password"
	"gorm.io/gorm"
)

// GORM 模型钩子 是在创建、查询、更新、删除等操作之前、之后调用的函数。
// 为模型定义指定的方法，它会在创建、更新、查询、删除时自动被调用。如果任何回调返回错误，GORM 将停止后续的操作并回滚事务
// Gorm的模型钩子 在模型创建前调用
// func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	user.Password = password.Hash(user.Password)
// 	return
// }

// // Gorm的模型钩子 在模型更新前调用
// func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
// 	if !password.IsHashed(user.Password) {
// 		user.Password = password.Hash(user.Password)
// 	}
// 	return
// }

// BeforeSave GORM 的模型钩子，在保存和更新模型前调用
func (u *User) BeforeSave(tx *gorm.DB) (err error) {

	if !password.IsHashed(u.Password) {
		u.Password = password.Hash(u.Password)
	}
	return
}
