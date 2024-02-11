package models

import "time"

type User struct {
	IdUser         string    `gorm:"primaryKey" json:"id_user"`
	Name           string    `json:"name"`
	Level          string    `json:"level"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	PasswordString string    `json:"password_string"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
