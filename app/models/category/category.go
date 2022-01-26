package category

import (
	"github.com/StubbornYouth/goblog/app/models"
	"github.com/StubbornYouth/goblog/pkg/route"
)

type Category struct {
	models.BaseModel
	Name string `gorm:"type:varchar(255);not null;" valid:"name"`
}

func (c Category) Link() string {
	return route.RouteNameToURL("categories.show", "id", c.GetStringID())
}
