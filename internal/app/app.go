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
	sessionRepo := postgres.NewPostgresSessionRepository(db, rdb)
	authService := service.NewAuthService(userRepo, sessionRepo, rdb)

	return &HandlersContainer{
		AuthHandler: handler.NewAuthHandler(authService, rdb),
	}
}
