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

	r.Post("/v1/auth/login", handlersContainer.AuthHandler.Login)

	r.Post("/v1/auth/forgot-password", handlersContainer.AuthHandler.ForgotPassword)

	r.Post("/v1/auth/verify-otp", handlersContainer.AuthHandler.VerifyOtp)

	r.Post("/v1/auth/reset-password", handlersContainer.AuthHandler.ResetPassword)

	return r
}
