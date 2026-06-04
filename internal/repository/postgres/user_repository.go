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

func (r *PostgresUserRepository) GetUserByUserName(userName_ string) (model.User, error) {
	var user model.User
	query := `
		SELECT name, email, mobile FROM users WHERE user_name = $1;
	`
	err := r.db.QueryRow(query, userName_).Scan(&user.Name, &user.Email, &user.Mobile)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	query := `
		SELECT user_id, name, mobile FROM users WHERE email = $1;
	`
	err := r.db.QueryRow(query, email).Scan(&user.UserID, &user.Name, &user.Mobile)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUserCredentials(userName_ string) (string, string, error) {
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

func (r *PostgresUserRepository) UpdatePassword(userName_, password_ string) error {
	query := `
		UPDATE passwords SET password = $1 WHERE user_name = $2;
	`
	_, err := r.db.Exec(query, password_, userName_)
	if err != nil {
		return err
	}
	return nil
}
