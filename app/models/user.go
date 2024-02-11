package models

import (
	"time"
)

type Level int

const (
	client Level = iota
	admin
	owner
)

// Define a method to return the string representation of the enum value
func (s Level) String() string {
	statuses := [...]string{"client", "admin", "owner"}
	if s < 0 || s >= Level(len(statuses)) {
		return "Unknown"
	}
	return statuses[s]
}

type User struct {
	IdUser         uint        `gorm:"primaryKey" json:"id_user"`
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	PasswordString string      `json:"password_string"`
	Level          interface{} `json:"level"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
