package router

import (
	"login/internal/database"
	"login/internal/handler"
	"login/internal/repository/postgres"
	"login/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitAuthModule() *handler.Handler {
	// Any necessary initialization can be done here
	db, err := database.GetDBInstance()
	if err != nil {
		panic("Failed to get database instance: " + err.Error())
	}
	userRepo := postgres.NewPostgresUserRepository(db)
	authService := service.NewAuthService(userRepo)
	return handler.NewHandler(authService)

}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Testing Router!!!"))
	})

	r.Post("/v1/auth/signup", InitAuthModule().SignupHandler)

	return r
}
