package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserID    uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string     `json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `gorm:"default:now()"`
	UpdatedAt *time.Time `gorm:"default:now()"`
	UserTOTP  UserTOTP   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserTOTP struct {
	UserTOTPID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId     uuid.UUID `gorm:"unique;type:uuid"`
	Secret     string
	IsEnabled  bool
	LastUsedAt *time.Time
	CreatedAt  *time.Time `gorm:"default:now()"`
}

type UserSessions struct {
	gorm.Model
}
