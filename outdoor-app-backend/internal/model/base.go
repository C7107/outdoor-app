package model

import (
	"time"
	"gorm.io/gorm"
)

// BaseModel 所有数据表的基础结构
type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除标识，不在前端返回
}