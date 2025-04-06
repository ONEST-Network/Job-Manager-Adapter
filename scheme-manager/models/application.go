package models

import (
	"time"
	"gorm.io/gorm"
)

type Application struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SchemeID    uint           `json:"scheme_id"`
	UserID      string         `json:"user_id"`      // unique user identifier
	UserName    string         `json:"user_name"`    // optional
	Credentials string         `json:"credentials"`  // credentials to validate eligibility (JSON format)

	Status      string         `json:"status"`       // e.g., "pending", "approved", "rejected"
	Remarks     string         `json:"remarks"`      // optional

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}