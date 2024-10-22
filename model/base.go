package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id         int64          `gorm:"column:id;primary_key" json:"id"`
	CreateTime time.Time      `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time      `gorm:"column:update_time" json:"update_time"`
	DeleteTime gorm.DeletedAt `gorm:"column:delete_time" json:"delete_time"`
}

// gorm 钩子函数，用于创建时间和更新时间的自动赋值
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreateTime = time.Now()
	b.UpdateTime = time.Now()
	return nil
}
