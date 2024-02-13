package models

import (
	"time"
)

type ResetCodePassword struct {
	Email     string    `gorm:"primaryKey" json:"email"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
