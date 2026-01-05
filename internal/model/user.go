package model

type User struct {
	NickName string `gorm:"column:nick_name;type:varchar(255);not null"`
	Email    string `gorm:"column:email;type:varchar(255);unique"`
	Phone    string `gorm:"column:phone;type:varchar(20);not null;unique"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
	CommonPartUnique
}
