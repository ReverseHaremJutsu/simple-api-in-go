package repository

import (
	"database/sql"

	"rest-api-in-gin/internal/account/domain/model/entity"
)

// PostgresUserAccountWriter is an implementation of domain interface, UserAccountWriter
type PostgresUserAccountWriter struct {
	db *sql.DB
}

// NewUserAccountWriter creates a new instance of UserAccountWriter
func NewPostgresUserAccountWriter(db *sql.DB) *PostgresUserAccountWriter {
	return &PostgresUserAccountWriter{db: db}
}

// Create persists the UserAccount into `account_user` table
func (r *PostgresUserAccountWriter) Create(user *entity.UserAccount) error {
	_, err := r.db.Exec(
		"INSERT INTO account_user (id, name, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.ID, user.Name.Value(), user.Email.Value(), user.Password.Value(), user.Role.Value(), user.CreatedAt, user.UpdatedAt,
	)
	return err
}
