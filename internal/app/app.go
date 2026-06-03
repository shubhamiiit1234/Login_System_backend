package app

import (
	"database/sql"
	"login/internal/handler"
	"login/internal/repository/postgres"
	"login/internal/service"

	"github.com/go-redis/redis/v8"
)

type HandlersContainer struct {
	AuthHandler *handler.AuthHandler
}

func NewHandlersContainer(db *sql.DB, rdb *redis.Client) *HandlersContainer {

	userRepo := postgres.NewPostgresUserRepository(db, rdb)
	authService := service.NewAuthService(userRepo, rdb)

	return &HandlersContainer{
		AuthHandler: handler.NewAuthHandler(authService),
	}
}
