package user

import (
	"github.com/StubbornYouth/goblog/app/models"
	"github.com/StubbornYouth/goblog/pkg/password"
	"github.com/StubbornYouth/goblog/pkg/route"
)

// 通过设置 GORM 模型的 Struct Tag 来解决字段数据库类型等设置
// GORM 默认会将键小写化作为字段名称，column 项可去除，另外默认是允许 NULL 的，故 default:NULL 项也可去除
// type User struct {
// 	models.BaseModel
// 	Name     string `gorm:"column:name;type:varchar(255);not null;unique;"`
// 	Email    string `gorm:"column:email;type:varchar(255);default:null;unique;"`
// 	Password string `gorm:"column:password;type:varchar(255);"`
// }
type User struct {
	models.BaseModel
	Name     string `gorm:"type:varchar(255);not null;unique;" valid:"name"`
	Email    string `gorm:"type:varchar(255);default:null;unique;" valid:"email"`
	Password string `gorm:"type:varchar(255);" valid:"password"`
	// gorm:"-" —— 设置 GORM 在读写时略过此字段，仅用于表单验证
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

// 匹配密码 此时存储的是明文 直接匹配
// func (user *User) ComparePassword(password string) bool {
// 	return user.Password == password
// }
func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}

// 生成用户链接
func (user User) Link() string {
	return route.RouteNameToURL("users.show", "id", user.GetStringID())
}
