package models

import "time"

type User struct {
	IdUser         string    `gorm:"primaryKey" json:"id_user"`
	Name           string    `json:"name"`
	Level          string    `json:"level"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	PasswordString string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
