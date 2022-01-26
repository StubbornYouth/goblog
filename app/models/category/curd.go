package category

import (
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/model"
	"github.com/StubbornYouth/goblog/pkg/types"
)

// 创建分类
func (category *Category) Create() (err error) {

	// article.ID             // 返回插入数据的主键
	// result.Error           // Create结果返回 error
	// result.RowsAffected    // 返回插入记录的条数
	if err := model.DB.Create(&category).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// 获取所有分类
func All() ([]Category, error) {
	var categories []Category

	if err := model.DB.Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

// 通过id获取分类
func Get(idstr string) (Category, error) {
	var category Category
	id := types.StringToInt(idstr)

	if err := model.DB.First(&category, id).Error; err != nil {
		return category, err
	}

	return category, nil
}
