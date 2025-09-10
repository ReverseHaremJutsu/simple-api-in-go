package repository

import (
	"database/sql"

	"rest-api-in-gin/internal/account/domain/model/value_object"
)

// PostgresUserAccountChecker is an implementation of domain interface, UserAccountChecker
type PostgresUserAccountChecker struct {
	db *sql.DB
}

// NewPostgresUserAccountChecker creates a new instance of PostgresUserAccountChecker
func NewPostgresUserAccountChecker(db *sql.DB) *PostgresUserAccountChecker {
	return &PostgresUserAccountChecker{db: db}
}

// ExistByEmail checks if given Email VO already exists inside `account_user` table
func (r *PostgresUserAccountChecker) ExistByEmail(email value_object.Email) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM account_user WHERE email = $1)`

	err := r.db.QueryRow(query, email.Value()).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
