package user

import (
	"web-chat/internal/model"

	"gorm.io/gorm"
)

type userDaoImpl struct {
	db *gorm.DB
}

func NewDao(db *gorm.DB) Dao {
	return &userDaoImpl{db: db}
}
func (u *userDaoImpl) CreateUser(user model.User) error {
	return u.db.Create(&user).Error
}

func (u *userDaoImpl) UpdateUser(uuid string, updateMap map[string]interface{}) error {
	return u.db.Model(&model.User{}).Where("uuid = ?", uuid).Updates(updateMap).Error
}

func (u *userDaoImpl) GetUserByID(uuid string) (*model.User, error) {
	var entity model.User
	if err := u.db.Where("uuid = ?", uuid).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (u *userDaoImpl) GetUserByPhone(phone string) (*model.User, error) {
	var entity model.User
	if err := u.db.Where("phone = ?", phone).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (u *userDaoImpl) GetUserByEmail(email string) (*model.User, error) {
	var entity model.User
	if err := u.db.Where("email = ?", email).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
