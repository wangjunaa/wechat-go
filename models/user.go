package models

import (
	"fmt"
	"time"
)

// ShowUser 可给他人展示的信息
type ShowUser struct {
	ID        string     `json:"id"`
	UserName  string     `json:"userName" gorm:"unique"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Birthday  *time.Time `json:"birthday"`
	LoginTime *time.Time `json:"loginTime"`
	Icon      []byte     `json:"icon"`
}

type UserBasic struct {
	ID          string `json:"id" `
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	UserName    string       `json:"userName" gorm:"unique"`
	Password    string       `json:"password" gorm:"not null"`
	Phone       string       `json:"phone"`
	Email       string       `json:"email"`
	HomeAddress string       `json:"homeAddress"`
	IpAddress   string       `json:"ipAddress"`
	LoginTime   *time.Time   `json:"loginTime"`
	Birthday    *time.Time   `json:"birthday"`
	Icon        []byte       `json:"icon"`
	Groups      []GroupBasic `json:"groups" gorm:"many2many:user_groups"`
}

func (user *UserBasic) TableName() string {
	return "Users"
}
func (user *UserBasic) Print() {
	fmt.Println(user.UserName, user.ID, user.Phone, user.Email)
}
func (user *UserBasic) ToShowUser() ShowUser {
	return ShowUser{
		ID:        user.ID,
		UserName:  user.UserName,
		Phone:     user.Phone,
		Email:     user.Email,
		Birthday:  user.Birthday,
		LoginTime: user.LoginTime,
		Icon:      user.Icon,
	}
}
