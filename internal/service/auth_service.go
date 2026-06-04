package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"login/internal/model"
	"login/internal/repository"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wneessen/go-mail"
)

var (
	ErrInvalidOTP      = errors.New("invalid OTP")
	ErrOTPExpired      = errors.New("OTP expired or not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrTooManyAttempts = errors.New("too many OTP verification attempts")

	TwoFAVerificationSuccessMessage = "2FA verification successful "
	OTPSent                         = "OTP has been sent to your email "
	OtpVerified                     = "OTP verification successful"
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	rdb         *redis.Client
}

// Dependency injection for AuthService
func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, rdb *redis.Client) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, rdb: rdb}
}

func (s *AuthService) Signup(user *model.User, password string) (int, error) {
	hashedPassword := hashPassword(password)
	return s.userRepo.CreateUser(user, hashedPassword)
}

func (s *AuthService) Login(userName_, password_ string) (string, string, error) {
	_, password, err := s.userRepo.GetUserCredentials(userName_)
	if err != nil {
		return "", "", err
	}

	if password != hashPassword(password_) {
		return "", "", ErrInvalidPassword
	}

	// Implementing 2FA with OTP
	otp := strconv.Itoa(rand.Intn(9000) + 1000)
	purpose := "2fa"

	userDetails, err := s.userRepo.GetUserByUserName(userName_)
	if err != nil {
		return "", "", err
	}

	redisKey := fmt.Sprintf("otp:%s:%s", purpose, userDetails.Email)

	err = s.rdb.HSet(context.Background(), redisKey, map[string]interface{}{
		"otp":      otp,
		"attempts": 0,
	}).Err()
	if err != nil {
		return "", "", err
	}

	s.rdb.Expire(context.Background(), redisKey, 10*time.Minute)
	go s.SendOtpOnEmail(userDetails.Email, otp)

	// Dummy session id
	session_id := "temp_session_" + strconv.Itoa(rand.Intn(1000000))
	return session_id, userDetails.Email, nil
}

func (s *AuthService) Complete2FA(ctx context.Context, purpose, email, otp string) (string, string, error) {
	message, err := s.VerifyOtp(purpose, email, otp)
	if err != nil {
		return "", "", err
	}

	userDetails, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	session_id := "final_session_" + strconv.Itoa(rand.Intn(1000000))

	err = s.sessionRepo.CreateSession(session_id, userDetails.UserID, email)
	if err != nil {
		return "", "", err
	}

	redisKey := fmt.Sprintf("session:%s", session_id)
	err = s.rdb.Set(context.Background(), redisKey, email, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting session in Redis:", err)
		return "", "", err
	}

	return session_id, message, nil
}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	otp := strconv.Itoa(rand.Intn(9000) + 1000)
	purpose := "forgot_password"
	redisKey := fmt.Sprintf("otp:%s:%s", purpose, email)

	err := s.rdb.HSet(context.Background(), redisKey, map[string]interface{}{
		"otp":      otp,
		"attempts": 0,
	}).Err()
	if err != nil {
		return "", err
	}

	s.rdb.Expire(context.Background(), redisKey, 10*time.Minute)

	go s.SendOtpOnEmail(email, otp)

	return OTPSent, nil
}

func (s *AuthService) VerifyOtp(purpose, email, otp string) (string, error) {
	redisKey := fmt.Sprintf("otp:%s:%s", purpose, email)
	data, err := s.rdb.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrOTPExpired
		}
		return "", err
	}

	if len(data) == 0 {
		return "", ErrOTPExpired
	}

	storedOtp := data["otp"]
	attempts, err := strconv.Atoi(data["attempts"])
	if err != nil {
		return "", err
	}

	if attempts >= 3 {
		s.rdb.Del(context.Background(), redisKey)
		return "", ErrTooManyAttempts
	}

	if storedOtp != otp {
		s.rdb.HIncrBy(context.Background(), redisKey, "attempts", 1)
		return "", ErrInvalidOTP
	}

	s.rdb.Del(context.Background(), redisKey)

	if purpose == "2fa" {
		return TwoFAVerificationSuccessMessage, nil
	}

	return OtpVerified, nil
}

func (s *AuthService) ResetPassword(userName, newPassword string) (string, error) {
	newHashedPassword := hashPassword(newPassword)
	err := s.userRepo.UpdatePassword(userName, newHashedPassword)
	if err != nil {
		return "FAILED TO RESET PASSWORD", err
	}

	return "SUCCESSFULLY RESET THE PASSWORD", nil
}

func (s *AuthService) SendOtpOnEmail(email, otp string) error {
	fromEmail := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	fmt.Println("fromEmail: ", fromEmail)
	fmt.Println("password: ", password)
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
		mail.WithPassword(password),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return err
	}

	if err := client.DialAndSend(m); err != nil {
		log.Printf("Failed to send: %v", err)
		return err
	}

	return nil
}

func hashPassword(password string) string {
	return password + "_hashed"
}
