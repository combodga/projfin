package store

import (
	"context"
	"fmt"

	"github.com/combodga/projfin"
	"github.com/jmoiron/sqlx"
)

type WithdrawPG struct {
	DB *sqlx.DB
}

func NewWithdrawPG(db *sqlx.DB) *WithdrawPG {
	return &WithdrawPG{DB: db}
}

func (w *WithdrawPG) ListWithdrawals(ctx context.Context, username string) ([]projfin.Withdraw, error) {
	var result []projfin.Withdraw

	withdraw1 := projfin.Withdraw{}
	sql := "SELECT * FROM withdrawals WHERE username = $1"
	rows, err := w.DB.QueryxContext(ctx, sql, username)
	if err != nil {
		return result, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&withdraw1)
		if err != nil {
			return result, fmt.Errorf("store scan error: %w", err)
		}
		result = append(result, withdraw1)
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("store get rows error: %w", err)
	}

	return result, nil
}

func (w *WithdrawPG) Withdraw(ctx context.Context, username, orderNumber string, sum float64) (int, error) {
	tx, err := w.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("store tx error: %w", err)
	}

	var uUsername string
	var uPassword string
	var uBalance float64
	var uWithdrawn float64
	sql := "SELECT * FROM users WHERE username = $1"
	err = tx.QueryRowContext(ctx, sql, username).Scan(&uUsername, &uPassword, &uBalance, &uWithdrawn)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("store query rows error: %w", err)
	}
	if uBalance < sum {
		return 402, nil
	}

	sql = "UPDATE users SET balance = $1, withdrawn = $2 WHERE username = $3"
	_, err = tx.ExecContext(ctx, sql, uBalance-sum, uWithdrawn+sum, username)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("store query error: %w", err)
	}

	sql = "INSERT INTO withdrawals VALUES ($1, $2, $3, NOW())"
	_, err = tx.ExecContext(ctx, sql, orderNumber, username, sum)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("store query error: %w", err)
	}

	return 0, tx.Commit()
}
