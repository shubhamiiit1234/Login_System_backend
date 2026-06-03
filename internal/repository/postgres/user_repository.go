package postgres

import (
	"database/sql"
	"login/internal/model"

	"github.com/go-redis/redis/v8"
)

type PostgresUserRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

// Dependency injection for PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB, rdb *redis.Client) *PostgresUserRepository {
	return &PostgresUserRepository{db: db, rdb: rdb}
}

type UserRepository interface {
	CreateUser(user *model.User, hashedPassword string) error
}

func (r *PostgresUserRepository) CreateUser(user *model.User, hashedPassword string) (int, error) {
	txn, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO users (name, email, mobile, user_name)
		VALUES ($1, $2, $3, $4) returning user_id;
	`
	var userID int
	err = r.db.QueryRow(query, user.Name, user.Email, user.Mobile, user.UserName).Scan(&userID)
	if err != nil {
		txn.Rollback()
		return 0, err
	}

	query = `
		INSERT INTO passwords (user_id, password, user_name)
		VALUES ($1, $2, $3);
	`
	_, err = r.db.Exec(query, userID, hashedPassword, user.UserName)

	if err != nil {
		txn.Rollback()
		return 0, err
	}

	err = txn.Commit()
	if err != nil {
		return 0, err
	}

	return userID, err
}

func (r *PostgresUserRepository) GetUserCredentials(userName_, password_ string) (string, string, error) {
	var userName, password string
	query := `
		SELECT user_name, password FROM passwords WHERE user_name = $1;
	`
	err := r.db.QueryRow(query, userName_).Scan(&userName, &password)
	if err != nil {
		return "", "", err
	}

	return userName, password, nil
}
