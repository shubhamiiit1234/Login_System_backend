package handler

import (
	"encoding/json"
	"login/internal/model"
	"login/internal/service"
	"net/http"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SignupResponse struct {
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// Common
type Handler struct {
	service *service.AuthService
}

// Dependency injection for Handler
func NewHandler(service *service.AuthService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signup Handler!!!"))

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		UserName: req.UserName,
	}

	err := h.service.Signup(user, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
