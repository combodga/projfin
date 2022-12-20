package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Store struct {
	DB        *sqlx.DB
	ErrorDupe error
}

type user struct {
	Username  string  `db:"username"`
	Password  string  `db:"password"`
	Balance   float64 `db:"balance"`
	Withdrawn float64 `db:"withdrawn"`
}
type order struct {
	OrderNumber string  `db:"order_number"`
	Username    string  `db:"username"`
	Status      string  `db:"status"`
	Accrual     float64 `db:"accrual"`
	UploadedAt  string  `db:"uploaded_at"`
}
type withdraw struct {
	OrderNumber string  `db:"order_number"`
	Username    string  `db:"username"`
	Sum         float64 `db:"sum"`
	ProcessedAt string  `db:"processed_at"`
}

func New(database string) (*Store, error) {
	db, err := sqlx.Connect("postgres", database)
	if err != nil {
		return nil, fmt.Errorf("store connect error: %w", err)
	}

	s := &Store{
		DB:        db,
		ErrorDupe: fmt.Errorf("duplicate key error"),
	}

	sql := "BEGIN;CREATE TABLE IF NOT EXISTS users (username text primary key,password text,balance double precision,withdrawn double precision);"
	sql += "CREATE TABLE IF NOT EXISTS orders (order_number text primary key,username text,status text,accrual double precision,uploaded_at timestamp with time zone);"
	sql += "CREATE TABLE IF NOT EXISTS withdrawals (order_number text,username text,sum double precision,processed_at timestamp with time zone);COMMIT;"
	db.MustExec(sql)

	user1 := user{}
	sql = "SELECT * FROM users"
	rows, err := db.Queryx(sql)
	if err != nil {
		return s, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.StructScan(&user1)
		if err != nil {
			return s, fmt.Errorf("store scan error: %w", err)
		}
	}
	err = rows.Err()
	if err != nil {
		return s, fmt.Errorf("store get rows error: %w", err)
	}

	return s, nil
}

func (s *Store) DoRegister(username, password string) error {
	sql := "INSERT INTO users VALUES ($1, $2, 0, 0)"
	_, err := s.DB.Exec(sql, username, password)
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

func (s *Store) DoLogin(username, password string) error {
	var count int
	sql := "SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2"
	rows, err := s.DB.Queryx(sql, username, password)
	if err != nil {
		return fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return fmt.Errorf("store query rows error: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("auth error: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("store get rows error: %w", err)
	}

	return nil
}

func (s *Store) CheckOrder(username, orderNumber string) (int, error) {
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

func (s *Store) MakeOrder(username, orderNumber string) error {
	sql := "INSERT INTO orders VALUES ($1, $2, 'NEW', 0, NOW())"
	_, err := s.DB.Exec(sql, orderNumber, username)
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

func (s *Store) ListOrders(username string) ([]order, error) {
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

func (s *Store) ListWithdrawals(username string) ([]withdraw, error) {
	var result []withdraw

	withdraw1 := withdraw{}
	sql := "SELECT * FROM withdrawals WHERE username = $1"
	rows, err := s.DB.Queryx(sql, username)
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

func (s *Store) GetUserBalance(username string) (user, error) {
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

func (s *Store) Withdraw(username, orderNumber string, sum float64) (int, error) {
	user1 := user{}
	sql := "SELECT * FROM users WHERE username = $1"
	rows, err := s.DB.Queryx(sql, username)
	if err != nil {
		return 0, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(&user1); err != nil {
			return 0, fmt.Errorf("store query rows error: %w", err)
		}
		if user1.Balance < sum {
			return 402, nil
		}
	}
	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("store get rows error: %w", err)
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("store tx error: %w", err)
	}

	sql = "UPDATE users SET balance = $1, withdrawn = $2 WHERE username = $3"
	_, err = tx.Exec(sql, user1.Balance-sum, user1.Withdrawn+sum, username)
	if err != nil {
		return 0, fmt.Errorf("store query error: %w", err)
	}

	sql = "INSERT INTO withdrawals VALUES ($1, $2, $3, NOW())"
	_, err = tx.Exec(sql, orderNumber, username, sum)
	if err != nil {
		return 0, fmt.Errorf("store query error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("store commit error: %w", err)
	}

	return 0, nil
}

func (s *Store) InvalidateOrder(orderNumber string) error {
	sql := "UPDATE orders SET status = 'INVALID' WHERE order_number = $1"
	_, err := s.DB.Exec(sql, orderNumber)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}

	return nil
}

func (s *Store) GetOrdersUser(orderNumber string) (order, error) {
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

func (s *Store) ProcessOrder(orderNumber string, accrual float64) error {
	order, err := s.GetOrdersUser(orderNumber)
	if err != nil {
		return fmt.Errorf("store user balance error: %w", err)
	}
	username := order.Username

	balance, err := s.GetUserBalance(username)
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

func (s *Store) OrdersProcessing() ([]order, error) {
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
