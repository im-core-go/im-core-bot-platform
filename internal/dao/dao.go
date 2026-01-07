package dao

import (
	"web-chat/internal/dao/chat"
	"web-chat/internal/dao/user"

	"gorm.io/gorm"
)

type Dao struct {
	UserDao user.Dao
	ChatDao chat.Dao
}

func NewDao(db *gorm.DB) *Dao {
	return &Dao{
		UserDao: user.NewDao(db),
		ChatDao: chat.NewDao(db),
	}
}
