package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"rest-api-in-gin/internal/payment/domain/model/entity"
	"rest-api-in-gin/internal/payment/domain/model/value_object"

	"github.com/google/uuid"
)

// PostgresWalletReader is an implementation of domain interface, WalletReader
type PostgresWalletReader struct {
	db *sql.DB
}

// NewPostgresWalletReader creates a new instance of PostgresWalletReader
func NewPostgresWalletReader(db *sql.DB) *PostgresWalletReader {
	return &PostgresWalletReader{db: db}
}

// GetByID retrieves a wallet by ID
func (r *PostgresWalletReader) GetByID(id uuid.UUID) (*entity.Wallet, error) {
	query := `SELECT id, user_account_id, balance FROM wallets WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var walletID uuid.UUID
	var userAccountID uuid.UUID
	var balance uint64

	err := row.Scan(&walletID, &userAccountID, &balance)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByID query scan error: %w", err)
	}

	moneyVO, err := value_object.NewMoney(balance)
	if err != nil {
		return nil, fmt.Errorf("invalid money value from db: %w", err)
	}

	return &entity.Wallet{
		ID:            walletID,
		UserAccountID: userAccountID,
		Balance:       moneyVO,
	}, nil
}
