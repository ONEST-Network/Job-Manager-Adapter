package models

import (
	"time"
	"gorm.io/gorm"
)

type Scheme struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Eligibility   string         `json:"eligibility"` // JSON or string rule, e.g., "age > 18"
	StartDate     time.Time      `json:"start_date"`
	EndDate       time.Time      `json:"end_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}