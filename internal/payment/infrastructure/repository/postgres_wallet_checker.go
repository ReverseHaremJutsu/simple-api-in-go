package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

// PostgresWalletChecker is an implementation of domain interface, WalletChecker
type PostgresWalletChecker struct {
	db *sql.DB
}

// NewPostgresWalletChecker creates a new instance of PostgresWalletChecker
func NewPostgresWalletChecker(db *sql.DB) *PostgresWalletChecker {
	return &PostgresWalletChecker{db: db}
}

// ExistByID checks if a wallet exists in the database
func (r *PostgresWalletChecker) ExistByID(walletID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1)`
	err := r.db.QueryRow(query, walletID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
