package article

import (
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/model"
	"github.com/StubbornYouth/goblog/pkg/types"
)

// 通过id获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)

	// First() 是 gorm.DB 提供的用以从结果集中获取第一条数据的查询方法
	// .Error 是 GORM 的错误处理机制。与常见的 Go 代码不同，因 GORM 提供的是链式 API，如果遇到任何错误，GORM 会设置 *gorm.DB 的 Error 字段，您需要像这样检查它。
	// 在 GORM 中，当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// 获取文章列表
// 将 map 类型的 Article 对象传参到 Find() 方法内，即可获取到所有文章数据
// 除了代码精简、可读性以外，不需要时刻记住关闭连接，也是使用 GORM 的优势之一
func GetAll() ([]Article, error) {
	var articles []Article

	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}

	return articles, nil
}

// 创建文章
func (article *Article) Create() (err error) {

	// article.ID             // 返回插入数据的主键
	// result.Error           // Create结果返回 error
	// result.RowsAffected    // 返回插入记录的条数
	if err := model.DB.Create(&article).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// 更新文章
func (article *Article) Update() (RowsAffected int64, err error) {
	result := model.DB.Save(&article)

	// result 返回两个参数 result.rowsAffected更新变化条数 result.Error 错误信息
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil

}

// 删除文章
func (article *Article) Delete() (RowsAffected int64, err error) {
	result := model.DB.Delete(&article)

	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}
