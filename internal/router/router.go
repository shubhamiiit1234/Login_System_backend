package router

import (
	"login/internal/app"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handlersContainer *app.HandlersContainer) http.Handler {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Testing Router!!!"))
	})

	r.Post("/v1/auth/signup", handlersContainer.AuthHandler.Signup)

	return r
}
