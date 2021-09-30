package article

import (
	"strconv"

	"github.com/StubbornYouth/goblog/pkg/route"
)

// article文章模型
type Article struct {
	ID    int64
	Title string
	Body  string
}

func (a Article) Link() string {
	return route.RouteNameToURL("articles.show", "id", strconv.FormatInt(a.ID, 10))
}
