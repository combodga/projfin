package withdrawstore

import (
	"context"
	"fmt"

	"github.com/combodga/projfin/internal/store"
)

type user struct {
	Username  string  `db:"username"`
	Password  string  `db:"password"`
	Balance   float64 `db:"balance"`
	Withdrawn float64 `db:"withdrawn"`
}

type withdraw struct {
	OrderNumber string  `db:"order_number"`
	Username    string  `db:"username"`
	Sum         float64 `db:"sum"`
	ProcessedAt string  `db:"processed_at"`
}

func ListWithdrawals(s *store.Store, ctx context.Context, username string) ([]withdraw, error) {
	var result []withdraw

	withdraw1 := withdraw{}
	sql := "SELECT * FROM withdrawals WHERE username = $1"
	rows, err := s.DB.QueryxContext(ctx, sql, username)
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

func GetUserBalance(s *store.Store, username string) (user, error) {
	user1 := user{}

	sql := "SELECT * FROM users WHERE username = $1"
	rows, err := s.DB.Queryx(sql, username)
	if err != nil {
		return user1, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(&user1)
	if err != nil {
		return user1, fmt.Errorf("store scan error: %w", err)
	}

	err = rows.Err()
	if err != nil {
		return user1, fmt.Errorf("store get rows error: %w", err)
	}

	return user1, nil
}

func Withdraw(s *store.Store, ctx context.Context, username, orderNumber string, sum float64) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("store tx error: %w", err)
	}

	user1 := user{}
	sql := "SELECT * FROM users WHERE username = $1"
	err = tx.QueryRowContext(ctx, sql, username).Scan(&user1)
	if err != nil {
		return 0, fmt.Errorf("store query rows error: %w", err)
	}
	if user1.Balance < sum {
		return 402, nil
	}

	sql = "UPDATE users SET balance = $1, withdrawn = $2 WHERE username = $3"
	_, err = tx.ExecContext(ctx, sql, user1.Balance-sum, user1.Withdrawn+sum, username)
	if err != nil {
		return 0, fmt.Errorf("store query error: %w", err)
	}

	sql = "INSERT INTO withdrawals VALUES ($1, $2, $3, NOW())"
	_, err = tx.ExecContext(ctx, sql, orderNumber, username, sum)
	if err != nil {
		return 0, fmt.Errorf("store query error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("store commit error: %w", err)
	}

	return 0, nil
}
