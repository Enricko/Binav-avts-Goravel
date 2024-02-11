package models

import "time"

type Client struct {
	IdClient  uint      `gorm:"primaryKey" json:"id_client"`
	User      User      `json:"id_user"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
