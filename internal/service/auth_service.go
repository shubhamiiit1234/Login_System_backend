package service

import (
	"context"
	"errors"
	"log"
	"login/internal/model"
	"login/internal/repository/postgres"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wneessen/go-mail"
)

type AuthService struct {
	userRepo *postgres.PostgresUserRepository
	rdb      *redis.Client
}

// Dependency injection for AuthService
func NewAuthService(userRepo *postgres.PostgresUserRepository, rdb *redis.Client) *AuthService {
	return &AuthService{userRepo: userRepo, rdb: rdb}
}

func (s *AuthService) Signup(user *model.User, password string) (int, error) {
	hashedPassword := hashPassword(password)
	return s.userRepo.CreateUser(user, hashedPassword)
}

func (s *AuthService) Login(userName_, password_ string) (string, error) {
	_, password, err := s.userRepo.GetUserCredentials(userName_, password_)
	if err != nil {
		return "", err
	}

	if password != hashPassword(password_) {
		return "", errors.New("Invalid Password")
	}

	// Dummy session id
	session_id := "sess_" + strconv.Itoa(rand.Intn(1000000))
	return session_id, nil

}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	otp := strconv.Itoa(rand.Intn(9000) + 1000)

	err := s.rdb.Set(context.Background(), "otp:"+email, otp, 10*time.Minute).Err()
	if err != nil {
		return "", err
	}

	go s.SendOtpOnEmail(email, otp)

	return "Password reset link has been sent to your email ", nil
}

func (s *AuthService) SendOtpOnEmail(email, otp string) {
	fromEmail := "gammaop3850@gmail.com"
	recipient := email

	m := mail.NewMsg()
	m.From(fromEmail)
	m.To(recipient)
	m.Subject("Testing login with otp")
	m.SetBodyString(mail.TypeTextPlain, "Your OTP is: "+otp)

	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(fromEmail),
		mail.WithPassword("Paste your Gmail app password here"),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}

	if err := client.DialAndSend(m); err != nil {
		log.Printf("Failed to send: %v", err)
		return
	}

}

func hashPassword(password string) string {
	return password + "_hashed"
}
