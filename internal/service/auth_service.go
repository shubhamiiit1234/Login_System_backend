package service

import (
	"fmt"
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
	fmt.Println("password", password)
	hashedPassword := hashPassword(password)
	fmt.Println("hashedPassword", hashedPassword)
	return s.userRepo.CreateUser(user, hashedPassword)
}

func hashPassword(password string) string {
	return password + strconv.Itoa(rand.Intn(100))
}
