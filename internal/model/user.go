package model

import "time"

type User struct {
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	UserName  string    `json:"user_name"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
