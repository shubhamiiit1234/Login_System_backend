package postgres

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

type PostgresSessionRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

// Dependency injection for PostgresSessionRepository
func NewPostgresSessionRepository(db *sql.DB, rdb *redis.Client) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db, rdb: rdb}
}

type SessionRepository interface {
	CreateSession(sessionID string, userID int, email string) error
	GetSession(sessionID string) (int, string, error)
}

func (r *PostgresSessionRepository) CreateSession(sessionID string, userID int, email string) error {
	query := `
		INSERT INTO sessions (session_id, user_id, email)
		VALUES ($1, $2, $3) RETURNING session_id;
	`
	_, err := r.db.Exec(query, sessionID, userID, email)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSessionRepository) GetSession(sessionID string) (int, string, error) {
	query := `
		SELECT user_id, email FROM sessions
		WHERE session_id = $1;
	`
	var userID int
	var email string
	err := r.db.QueryRow(query, sessionID).Scan(&userID, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", nil // Session not found
		}
		return 0, "", err
	}
	return userID, email, nil
}
