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

func (r *PostgresUserRepository) CreateUser(user *model.User, hashedPassword string) error {
	query := `
		INSERT INTO users (name, email, mobile, user_name, password, verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Mobile, user.UserName, hashedPassword, user.Verified, user.CreatedAt, user.UpdatedAt)
	return err
}
