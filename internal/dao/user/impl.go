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

func (u *userDaoImpl) UpdateUser(id int64, updateMap map[string]interface{}) error {
	return u.db.Model(&model.User{}).Where("id = ?", id).Updates(updateMap).Error
}

func (u *userDaoImpl) GetUserByID(id int64) (*model.User, error) {
	var entity model.User
	if err := u.db.First(&entity, id).Error; err != nil {
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
