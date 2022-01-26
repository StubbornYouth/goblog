package article

import (
	"github.com/StubbornYouth/goblog/app/models"
	"github.com/StubbornYouth/goblog/app/models/category"
	"github.com/StubbornYouth/goblog/app/models/user"
	"github.com/StubbornYouth/goblog/pkg/route"
)

// article文章模型
// type Article struct {
// 	ID    int64
// 	Title string
// 	Body  string
// }

type Article struct {
	models.BaseModel
	Title string `gorm:"type:varchar(255);index" valid:"title"`
	Body  string `gorm:"type:varchar(255);index" valid:"body"`
	// GORM 是遵循约定优先于配置的风格。约定规定，一对一的关联，是通过模型单数小写加下划线和 ID 如 user_id 来实现，它会去读下划线前面的单词的复数形式作为数据表（如此例的 users），取字段 id 为关联字段。
	// 在特殊的情况下，如果你无法按照约定，也可以自行配置关联关系，具体请见 Belongs To 文档
	//  CompanyRefer int
	//  User   user.User `gorm:"foreignKey:CompanyRefer"` // 使用 CompanyRefer 作为外键}
	UserID     uint64 `gorm:"not null;index"`
	User       user.User
	CategoryID uint64 `gorm:"not null;default:4;index"`
	Category   category.Category
}

func (a Article) Link() string {
	// return route.RouteNameToURL("articles.show", "id", strconv.FormatInt(a.ID, 10))
	return route.RouteNameToURL("articles.show", "id", a.GetStringID())
}

// 格式化创建时间
func (a Article) CreatedAtDate() string {
	return a.CreatedAt.Format("2006-01-02")
}
