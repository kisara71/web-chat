package user

import "web-chat/internal/model"

type Dao interface {
	CreateUser(user model.User) error
	UpdateUser(id int64, updateMap map[string]interface{}) error
	GetUserByID(id int64) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}
