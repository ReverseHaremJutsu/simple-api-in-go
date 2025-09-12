package repository

import (
	"database/sql"
	"fmt"
	"time"

	"rest-api-in-gin/internal/account/domain/model/entity"

	"github.com/google/uuid"
)

// PostgresUserAccountWriter is an implementation of domain interface, UserAccountWriter
type PostgresUserAccountWriter struct {
	db *sql.DB
}

// NewUserAccountWriter creates a new instance of UserAccountWriter
func NewPostgresUserAccountWriter(db *sql.DB) *PostgresUserAccountWriter {
	return &PostgresUserAccountWriter{db: db}
}

// Create persists the UserAccount into `account_user` table and `outbox` table in a transaction
func (r *PostgresUserAccountWriter) Create(user *entity.UserAccount) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	_, err = r.db.Exec(`
		INSERT INTO account_user (id, name, email, password, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.Name.Value(), user.Email.Value(), user.Password.Value(), user.Role.Value(), user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	eventPayload := fmt.Sprintf(`{"id": "%s"}`, user.ID)
	_, err = tx.Exec(`
		INSERT INTO outbox (id, aggregate_name, aggregate_id, event_name, payload, topic_name, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		uuid.New(), "UserAccount", user.ID, "UserCreated", eventPayload, "user-created", time.Now()) // technically topic does not need equal event name, but im keeping it trivial for now
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
