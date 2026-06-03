package model

type Password struct {
	Id        int    `json:"id"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
