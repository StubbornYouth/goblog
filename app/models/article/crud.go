package article

import (
	"net/http"

	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/model"
	"github.com/StubbornYouth/goblog/pkg/pagination"
	"github.com/StubbornYouth/goblog/pkg/route"
	"github.com/StubbornYouth/goblog/pkg/types"
)

// 通过id获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)

	// First() 是 gorm.DB 提供的用以从结果集中获取第一条数据的查询方法
	// .Error 是 GORM 的错误处理机制。与常见的 Go 代码不同，因 GORM 提供的是链式 API，如果遇到任何错误，GORM 会设置 *gorm.DB 的 Error 字段，您需要像这样检查它。
	// 在 GORM 中，当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误
	// Preload 预加载用户信息
	if err := model.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// 获取文章列表
// 将 map 类型的 Article 对象传参到 Find() 方法内，即可获取到所有文章数据
// 除了代码精简、可读性以外，不需要时刻记住关闭连接，也是使用 GORM 的优势之一
// func GetAll() ([]Article, error) {
// 	var articles []Article
// 	// Debug 在终端显示执行的sql语句
// 	// if err := model.DB.Debug().Preload("User").Find(&articles).Error; err != nil {
// 	if err := model.DB.Preload("User").Find(&articles).Error; err != nil {
// 		return articles, err
// 	}

// 	return articles, nil
// }

// 增加分页数据
func GetAll(r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	// 初始化分页实例
	db := model.DB.Model(Article{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.RouteNameToURL("articles.index"), perPage)

	// 获取实例数据
	viewData := _pager.Paging()

	// 获取数据库数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
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

// 获取用户文章
func GetByUserID(idstr string) ([]Article, error) {
	var articles []Article

	if err := model.DB.Where("user_id = ?", idstr).Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}

	return articles, nil
}

// 获取分类文章
func GetByCategoryID(idstr string, r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	// 1. 初始化分页实例
	db := model.DB.Model(Article{}).Where("category_id = ?", idstr).Order("created_at desc")
	_pager := pagination.New(r, db, route.RouteNameToURL("categories.show", "id", idstr), perPage)

	// 获取视图数据
	viewData := _pager.Paging()

	// 获取文章数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}
