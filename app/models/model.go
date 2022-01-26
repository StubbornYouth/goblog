package models

import (
	"time"

	"github.com/StubbornYouth/goblog/pkg/types"
)

// 声明 GORM 数据模型时，字段标签是可选的，GORM 支持以下：（注：名大小写不敏感，但建议使用 camelCase 风格）
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;not null"`

	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

func (model BaseModel) GetStringID() string {
	//  按长度分为：int8、int16、int32、int64 对应的无符号整型：uint8、uint16、uint32、uint64
	return types.Uint64ToString(model.ID)
}
