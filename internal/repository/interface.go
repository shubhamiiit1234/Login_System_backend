package repository

import "login/internal/model"

type UserRepository interface {
	CreateUser(user *model.User, hashedPassword string) (int, error)
	GetUserByUserName(userName string) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	GetUserCredentials(userName string) (string, string, error)
	UpdatePassword(userName, password string) error
}

type SessionRepository interface {
	CreateSession(sessionID string, userID int, email string) error
	GetSession(sessionID string) (int, string, error)
}
