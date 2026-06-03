package postgres

import (
	"database/sql"
	"login/internal/model"
)

type PostgresUserRepository struct {
	db *sql.DB
}

// Dependency injection for PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

type UserRepository interface {
	CreateUser(user *model.User, hashedPassword string) error
}

func (r *PostgresUserRepository) CreateUser(user *model.User, hashedPassword string) (int, error) {
	query := `
		INSERT INTO users (name, email, mobile, user_name)
		VALUES ($1, $2, $3, $4) returning user_id;
	`
	var userID int
	err := r.db.QueryRow(query, user.Name, user.Email, user.Mobile, user.UserName).Scan(&userID)
	return userID, err
}
