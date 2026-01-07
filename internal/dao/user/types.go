package user

import "web-chat/internal/model"

type Dao interface {
	CreateUser(user model.User) error
	UpdateUser(uuid string, updateMap map[string]interface{}) error
	GetUserByID(uuid string) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}
