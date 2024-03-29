package models

import "time"

type Client struct {
	IdClient  string    `gorm:"primaryKey" json:"id_client"`
	IdUser    string    `json:"id_user"`
	User      *User     `gorm:"foreignKey:IdUser;references:IdUser"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}