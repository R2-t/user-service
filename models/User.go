package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type User struct {
	UserID       uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string         `json:"name"`
	Email        string         `gorm:"uniqueIndex" json:"email"`
	Password     string         `json:"password"`
	CreatedAt    *time.Time     `gorm:"default:now()"`
	UpdatedAt    *time.Time     `gorm:"default:now()"`
	UserTOTP     UserTOTP       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserSessions []UserSessions `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserTOTP struct {
	UserTOTPID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid"`
	Secret     string
	IsEnabled  bool
	LastUsedAt *time.Time
	CreatedAt  *time.Time `gorm:"default:now()"`
}

type UserSessions struct {
	UserSessionsID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TokenID        string    `gorm:"uniqueIndex"`
	TokenHash      string
	DeviceInfo     datatypes.JSON
	IsRevoked      bool       `gorm:"default:false"`
	CreatedAt      *time.Time `gorm:"default:now()"`
	UserID         uuid.UUID
}

type DeviceInfo struct {
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
	DeviceType string `json:"device_type,omitempty"`
	Browser    string `json:"browser,omitempty"`
	OS         string `json:"os,omitempty"`
}
