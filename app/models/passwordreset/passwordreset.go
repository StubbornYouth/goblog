package passwordreset

import "github.com/StubbornYouth/goblog/app/models"

type PasswordReset struct {
	models.BaseModel
	Token string `gorm:"type:varchar(255);index"`
	Email string `gorm:"type:varchar(255);index" valid:"email"`
}
