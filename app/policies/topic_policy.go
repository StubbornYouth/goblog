package policies

import (
	"github.com/StubbornYouth/goblog/app/models/article"
	"github.com/StubbornYouth/goblog/pkg/auth"
)

// 是否允许修改
func CanModifyArticle(_article article.Article) bool {
	return auth.User().ID == _article.UserID
}
