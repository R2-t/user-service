package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name     string   `json:"name"`
	Email    string   `gorm:"unique" json:"email"`
	Password string   `json:"password"`
	UserTOTP UserTOTP `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserTOTP struct {
	gorm.Model
	UserId     uint `gorm:"unique"`
	Secret     string
	IsEnabled  bool
	LastUsedAt *time.Time
}

type UserSessions struct {
	gorm.Model
}
