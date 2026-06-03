package service

import (
	"encoding/json"
	"errors"
	"login/internal/model"
	"login/internal/repository/postgres"
	"math/rand"
	"strconv"
)

type AuthService struct {
	userRepo *postgres.PostgresUserRepository
}

// Dependency injection for AuthService
func NewAuthService(userRepo *postgres.PostgresUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Signup(user *model.User, password string) (int, error) {
	hashedPassword := hashPassword(password)
	return s.userRepo.CreateUser(user, hashedPassword)
}

func (s *AuthService) Login(userName_, password_ string) (json.Token, error) {
	_, password, err := s.userRepo.GetUserCredentials(userName_, password_)
	if err != nil {
		return "", err
	}

	if password != hashPassword(password_) {
		return "", errors.New("Invalid Password")
	}

	// Dummy JWT token
	token := "token_" + strconv.Itoa(rand.Intn(1000000))
	return json.Token(token), nil

}

func hashPassword(password string) string {
	return password + "_hashed"
}
