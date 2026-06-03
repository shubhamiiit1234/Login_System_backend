package app

import (
	"database/sql"
	"login/internal/handler"
	"login/internal/repository/postgres"
	"login/internal/service"
)

type HandlersContainer struct {
	AuthHandler *handler.AuthHandler
}

func NewHandlersContainer(db *sql.DB) *HandlersContainer {

	userRepo := postgres.NewPostgresUserRepository(db)
	authService := service.NewAuthService(userRepo)

	return &HandlersContainer{
		AuthHandler: handler.NewAuthHandler(authService),
	}
}
