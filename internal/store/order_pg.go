package store

import (
	"context"
	"fmt"

	"github.com/combodga/projfin"
	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

type OrderPG struct {
	DB *sqlx.DB
}

func NewOrderPG(db *sqlx.DB) *OrderPG {
	return &OrderPG{DB: db}
}

func (o *OrderPG) CheckOrder(username, orderNumber string) (projfin.OrderStatus, error) {
	order1 := projfin.Order{}
	sql := "SELECT * FROM orders WHERE order_number = $1"
	rows, err := o.DB.Queryx(sql, orderNumber)
	if err != nil {
		return 0, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(&order1); err != nil {
			return projfin.OrderStatusError, fmt.Errorf("store query rows error: %w", err)
		}
		if order1.Username == username {
			return projfin.OrderStatusExists, nil
		} else if order1.Username != "" {
			return projfin.OrderStatusOccupied, nil
		}
	}

	err = rows.Err()
	if err != nil {
		return projfin.OrderStatusError, fmt.Errorf("store get rows error: %w", err)
	}

	return projfin.OrderStatusOK, nil
}

func (o *OrderPG) MakeOrder(ctx context.Context, username, orderNumber string) error {
	sql := "INSERT INTO orders VALUES ($1, $2, 'NEW', 0, NOW())"
	_, err := o.DB.ExecContext(ctx, sql, orderNumber, username)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == PGDuplicateCode {
				return ErrorDupe
			}
		}
		return fmt.Errorf("store query error: %w", err)
	}
	return nil
}

func (o *OrderPG) ListOrders(username string) ([]projfin.Order, error) {
	var result []projfin.Order

	order1 := projfin.Order{}
	sql := "SELECT * FROM orders WHERE username = $1"
	rows, err := o.DB.Queryx(sql, username)
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

func (o *OrderPG) InvalidateOrder(orderNumber string) error {
	sql := "UPDATE orders SET status = 'INVALID' WHERE order_number = $1"
	_, err := o.DB.Exec(sql, orderNumber)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}

	return nil
}

func (o *OrderPG) GetOrdersUser(orderNumber string) (projfin.Order, error) {
	order1 := projfin.Order{}

	sql := "SELECT * FROM orders WHERE order_number = $1"
	rows, err := o.DB.Queryx(sql, orderNumber)
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

func (o *OrderPG) ProcessOrder(orderNumber string, accrual float64) error {
	order, err := o.GetOrdersUser(orderNumber)
	if err != nil {
		return fmt.Errorf("store user balance error: %w", err)
	}
	username := order.Username

	balance, err := o.GetUserBalance(username)
	if err != nil {
		return fmt.Errorf("store user balance error: %w", err)
	}

	tx, err := o.DB.Begin()
	if err != nil {
		return fmt.Errorf("store tx error: %w", err)
	}

	sql := "UPDATE orders SET status = 'PROCESSED', accrual = $1 WHERE order_number = $2"
	_, err = tx.Exec(sql, accrual, orderNumber)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("db update error: %w", err)
	}

	sql = "UPDATE users SET balance = $1 WHERE username = $2"
	_, err = tx.Exec(sql, balance.Balance+accrual, username)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("store query error: %w", err)
	}

	return tx.Commit()
}

func (o *OrderPG) OrdersProcessing() ([]projfin.Order, error) {
	var result []projfin.Order

	sql := "UPDATE orders SET status = 'PROCESSING' WHERE status = 'NEW'"
	_, err := o.DB.Exec(sql)
	if err != nil {
		return result, fmt.Errorf("db update error: %w", err)
	}

	order1 := projfin.Order{}
	sql = "SELECT * FROM orders WHERE status = 'PROCESSING'"
	rows, err := o.DB.Queryx(sql)
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

func (o *OrderPG) GetUserBalance(username string) (projfin.User, error) {
	user1 := projfin.User{}

	sql := "SELECT * FROM users WHERE username = $1"
	rows, err := o.DB.Queryx(sql, username)
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
