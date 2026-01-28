package dao

import (
	"github.com/im-core-go/im-core-bot-platform/internal/dao/chat"

	"gorm.io/gorm"
)

type Dao struct {
	ChatDao chat.Dao
}

func NewDao(db *gorm.DB) *Dao {
	return &Dao{
		ChatDao: chat.NewDao(db),
	}
}
