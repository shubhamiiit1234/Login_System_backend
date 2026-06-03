package handler

import (
	"encoding/json"
	"fmt"
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

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	json.Token `json:"token"`
}

// Common
type AuthHandler struct {
	service *service.AuthService
}

// Dependency injection for Handler
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Signup Handler!!!"))

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

	user_id, err := h.service.Signup(user, req.Password)
	if err != nil {
		fmt.Println("Error creating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := SignupResponse{
		UserID:    user_id,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Login Handler!!!"))

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(req.UserName, req.Password)
	if err != nil {
		fmt.Println("Error logging in:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
