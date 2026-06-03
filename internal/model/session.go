package model

type Session struct {
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}
