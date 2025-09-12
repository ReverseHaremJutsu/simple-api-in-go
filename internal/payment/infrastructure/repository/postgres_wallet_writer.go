package repository

import (
	"database/sql"
	"fmt"
	"time"

	"rest-api-in-gin/internal/payment/domain/model/entity"
	"rest-api-in-gin/internal/payment/domain/repository"
)

// PostgresWalletWriter is an implementation of domain interface, WalletWriter
type PostgresWalletWriter struct {
	db *sql.DB
}

// NewPostgresWalletWriter creates a new instance of PostgresWalletWriter
func NewPostgresWalletWriter(db *sql.DB) *PostgresWalletWriter {
	return &PostgresWalletWriter{db: db}
}

// UpdatesByID inserts Payment rows and updates Wallet balances in a single transaction
func (r *PostgresWalletWriter) UpdatesByID(updates []repository.WalletPaymentUpdate) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	for _, u := range updates {

		p := u.Payment
		paymentQuery := `
			INSERT INTO payments (id, source_wallet_id, destination_wallet_id, amount, action, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`
		_, err := tx.Exec(paymentQuery, p.ID, p.SourceWallet(), p.DestinationWallet(), p.Amount(), p.PaymentAction())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert payment row: %w", err)
		}

		w := u.Wallet
		walletQuery := `
			UPDATE wallets 
			SET balance = $1, modified_at = NOW()
			WHERE id = $2
		`
		_, err = tx.Exec(walletQuery, w.Balance.Value(), w.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update wallet balance: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Create persists the Wallet into `wallet` table
func (r *PostgresWalletWriter) Create(wallet *entity.Wallet) error {
	now := time.Now()
	_, err := r.db.Exec(
		"INSERT INTO wallets (id, user_account_id, balance,created_at, modified_at) VALUES ($1, $2, $3, $4, $5)",
		wallet.ID, wallet.UserAccountID, wallet.Balance.Value(), now, now,
	)
	return err
}
