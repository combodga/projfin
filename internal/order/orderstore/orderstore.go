package orderstore

import (
	"context"
	"fmt"

	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/withdraw/withdrawstore"

	"github.com/lib/pq"
)

type order struct {
	OrderNumber string  `db:"order_number"`
	Username    string  `db:"username"`
	Status      string  `db:"status"`
	Accrual     float64 `db:"accrual"`
	UploadedAt  string  `db:"uploaded_at"`
}

func CheckOrder(s *store.Store, username, orderNumber string) (int, error) {
	order1 := order{}
	sql := "SELECT * FROM orders WHERE order_number = $1"
	rows, err := s.DB.Queryx(sql, orderNumber)
	if err != nil {
		return 0, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(&order1); err != nil {
			return 0, fmt.Errorf("store query rows error: %w", err)
		}
		if order1.Username == username {
			return 1, nil
		} else if order1.Username != "" {
			return 2, nil
		}
	}

	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("store get rows error: %w", err)
	}

	return 0, nil
}

func MakeOrder(s *store.Store, ctx context.Context, username, orderNumber string) error {
	sql := "INSERT INTO orders VALUES ($1, $2, 'NEW', 0, NOW())"
	_, err := s.DB.ExecContext(ctx, sql, orderNumber, username)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return s.ErrorDupe
			}
		}
		return fmt.Errorf("store query error: %w", err)
	}
	return nil
}

func ListOrders(s *store.Store, username string) ([]order, error) {
	var result []order

	order1 := order{}
	sql := "SELECT * FROM orders WHERE username = $1"
	rows, err := s.DB.Queryx(sql, username)
	if err != nil {
		return result, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&order1)
		if err != nil {
			return result, fmt.Errorf("store scan error: %w", err)
		}
		result = append(result, order1)
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("store get rows error: %w", err)
	}

	return result, nil
}

func InvalidateOrder(s *store.Store, orderNumber string) error {
	sql := "UPDATE orders SET status = 'INVALID' WHERE order_number = $1"
	_, err := s.DB.Exec(sql, orderNumber)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}

	return nil
}

func GetOrdersUser(s *store.Store, orderNumber string) (order, error) {
	order1 := order{}

	sql := "SELECT * FROM orders WHERE order_number = $1"
	rows, err := s.DB.Queryx(sql, orderNumber)
	if err != nil {
		return order1, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(&order1)
	if err != nil {
		return order1, fmt.Errorf("store scan error: %w", err)
	}

	err = rows.Err()
	if err != nil {
		return order1, fmt.Errorf("store get rows error: %w", err)
	}

	return order1, nil
}

func ProcessOrder(s *store.Store, orderNumber string, accrual float64) error {
	order, err := GetOrdersUser(s, orderNumber)
	if err != nil {
		return fmt.Errorf("store user balance error: %w", err)
	}
	username := order.Username

	balance, err := withdrawstore.GetUserBalance(s, username)
	if err != nil {
		return fmt.Errorf("store user balance error: %w", err)
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("store tx error: %w", err)
	}

	sql := "UPDATE orders SET status = 'PROCESSED', accrual = $1 WHERE order_number = $2"
	_, err = tx.Exec(sql, accrual, orderNumber)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}

	sql = "UPDATE users SET balance = $1 WHERE username = $2"
	_, err = tx.Exec(sql, balance.Balance+accrual, username)
	if err != nil {
		return fmt.Errorf("store query error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("store commit error: %w", err)
	}

	return nil
}

func OrdersProcessing(s *store.Store) ([]order, error) {
	var result []order

	sql := "UPDATE orders SET status = 'PROCESSING' WHERE status = 'NEW'"
	_, err := s.DB.Exec(sql)
	if err != nil {
		return result, fmt.Errorf("db update error: %w", err)
	}

	order1 := order{}
	sql = "SELECT * FROM orders WHERE status = 'PROCESSING'"
	rows, err := s.DB.Queryx(sql)
	if err != nil {
		return result, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&order1)
		if err != nil {
			return result, fmt.Errorf("store scan error: %w", err)
		}
		result = append(result, order1)
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("store get rows error: %w", err)
	}

	return result, nil
}
