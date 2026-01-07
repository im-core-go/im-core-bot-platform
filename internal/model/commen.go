package model

import "gorm.io/plugin/soft_delete"

type CommonPartUnique struct {
	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_name,sort:desc"`
}
type CommonPartNoUnique struct {
	CreatedAt int64 `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt
}
