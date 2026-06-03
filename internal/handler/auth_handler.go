package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"login/internal/model"
	"login/internal/service"
	"net/http"

	"github.com/go-redis/redis/v8"
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
	Session_id string `json:"session_id"`
	Status     string `json:"status"`
	Email      string `json:"email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type VerifyOtpRequest struct {
	Email   string `json:"email"`
	Otp     string `json:"otp"`
	Purpose string `json:"purpose"` // "2fa" or "forgot_password"
}

type VerifyOtpResponse struct {
	Message    string `json:"message"`
	Session_id string `json:"session_id,omitempty"`
}

// Common
type AuthHandler struct {
	service *service.AuthService
	rdb     *redis.Client
}

// Dependency injection for Handler
func NewAuthHandler(service *service.AuthService, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{service: service, rdb: rdb}
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

	session_id, email, err := h.service.Login(req.UserName, req.Password)
	if err != nil {
		fmt.Println("Error logging in:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{
		Session_id: session_id,
		Status:     "Pending 2FA verification",
		Email:      email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Forgot Password Handler!!!"))

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := h.service.ForgotPassword(req.Email)
	if err != nil {
		fmt.Println("Error in forgot password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := ForgotPasswordResponse{
		Message: message + req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) VerifyOtp(w http.ResponseWriter, r *http.Request) {

	var req VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Purpose != "2fa" && req.Purpose != "forgot_password" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := VerifyOtpResponse{}

	if req.Purpose == "2fa" {

		session_id, message, err := h.service.Complete2FA(context.Background(), req.Purpose, req.Email, req.Otp)
		if err != nil {
			fmt.Println("Error in completing 2FA:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.Session_id = session_id
		resp.Message = message
	} else {
		message, err := h.service.VerifyOtp(req.Purpose, req.Email, req.Otp)
		if err != nil {
			fmt.Println("Error in verifying OTP:", err)
			if err == service.ErrInvalidOTP || err == service.ErrOTPExpired {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		resp.Message = message
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
