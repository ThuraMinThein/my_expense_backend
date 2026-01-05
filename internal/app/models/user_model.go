package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Profile      string    `json:"profile"`
	Username     string    `json:"username" gorm:"unique"`
	Email        string    `json:"email" binding:"required" gorm:"unique"`
	Password     string    `json:"-"`
	GoogleID     string    `json:"google_id" gorm:"unique"`
	AuthProvider string    `json:"auth_provider" gorm:"default:'local'"`
	UserToken    UserToken `json:"-" gorm:"foreignKey:UserId"`
}

type UserToken struct {
	UserId       uint   `json:"-" gorm:"primaryKey;autoIncrement:false"`
	AccessToken  string `json:"access_token" gorm:"-"`
	RefreshToken string `json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (UserToken) TableName() string {
	return "user_tokens"
}
