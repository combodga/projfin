package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Store struct {
	DB        string
	ErrorDupe error
}

type User struct {
	Username  string  `db:"username"`
	Password  string  `db:"password"`
	Balance   float64 `db:"balance"`
	Withdrawn float64 `db:"withdrawn"`
}

type Order struct {
	OrderNumber string  `db:"order_number"`
	Username    string  `db:"username"`
	Status      string  `db:"status"`
	Accrual     float64 `db:"accrual"`
	UploadedAt  string  `db:"uploaded_at"`
	Sum         float64 `db:"sum"`
	ProcessedAt string  `db:"processed_at"`
}

func New(database string) (*Store, error) {
	s := &Store{
		DB:        database,
		ErrorDupe: fmt.Errorf("duplicate key error"),
	}

	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return s, fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	db.MustExec(`
        CREATE TABLE IF NOT EXISTS users (
            username text primary key,
            password text,
            balance double precision,
            withdrawn double precision
        );

        CREATE TABLE IF NOT EXISTS orders (
            order_number text primary key,
            username text,
            status text,
            accrual double precision,
            uploaded_at timestamp with time zone,
            sum double precision,
            processed_at timestamp with time zone
        );
    `)

	user := User{}
	rows, err := db.Queryx("SELECT * FROM users")
	if err != nil {
		return s, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.StructScan(&user)
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
	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users VALUES ($1, $2, 0, 0)", username, password)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return s.ErrorDupe
			}
		}
	}

	if err != nil {
		return fmt.Errorf("store query error: %w", err)
	}
	return nil
}

func (s *Store) DoLogin(username, password string) error {
	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	var count int
	rows, err := db.Queryx("SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2", username, password)
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
	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return 0, fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	order := Order{}
	rows, err := db.Queryx("SELECT * FROM orders WHERE order_number = $1", orderNumber)
	if err != nil {
		return 0, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(&order); err != nil {
			return 0, fmt.Errorf("store query rows error: %w", err)
		}
		if order.Username == username {
			return 1, nil
		} else if order.Username != "" {
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
	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO orders VALUES ($1, $2, 'NEW', 0, NOW(), 0, NOW())", orderNumber, username)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return s.ErrorDupe
			}
		}
	}

	if err != nil {
		return fmt.Errorf("store query error: %w", err)
	}

	return nil
}

func (s *Store) ListOrders(username string) ([]Order, error) {
	var result []Order

	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return result, fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	order := Order{}
	rows, err := db.Queryx("SELECT * FROM orders WHERE username = $1", username)
	if err != nil {
		return result, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&order)
		if err != nil {
			return result, fmt.Errorf("store scan error: %w", err)
		}
		result = append(result, order)
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("store get rows error: %w", err)
	}

	return result, nil
}

func (s *Store) GetUserBalance(username string) (User, error) {
	user := User{}

	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return user, fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	rows, err := db.Queryx("SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return user, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(&user)
	if err != nil {
		return user, fmt.Errorf("store scan error: %w", err)
	}

	err = rows.Err()
	if err != nil {
		return user, fmt.Errorf("store get rows error: %w", err)
	}

	return user, nil
}

func (s *Store) Withdraw(username, orderNumber string, sum float64) (int, error) {
	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return 0, fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	user := User{}
	rows, err := db.Queryx("SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return 0, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return 0, fmt.Errorf("store query rows error: %w", err)
		}
		if user.Balance < sum {
			return 402, nil
		}
	}
	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("store get rows error: %w", err)
	}

	_, err = db.Exec("UPDATE users SET balance = $1, withdrawn = $2 WHERE username = $3", user.Balance-sum, user.Withdrawn+sum, username)
	if err != nil {
		return 0, fmt.Errorf("store query error: %w", err)
	}

	_, err = db.Exec("UPDATE orders SET sum = $1, processed_at = NOW() WHERE order_number = $2", sum, orderNumber)
	if err != nil {
		return 0, fmt.Errorf("store query error: %w", err)
	}

	return 0, nil
}

func (s *Store) InvalidateOrder(orderNumber string) error {
	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE orders SET status = 'INVALID' WHERE order_number = $1", orderNumber)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}

	return nil
}

func (s *Store) GetOrdersUser(orderNumber string) (Order, error) {
	order := Order{}

	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return order, fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	rows, err := db.Queryx("SELECT * FROM orders WHERE order_number = $1", orderNumber)
	if err != nil {
		return order, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(&order)
	if err != nil {
		return order, fmt.Errorf("store scan error: %w", err)
	}

	err = rows.Err()
	if err != nil {
		return order, fmt.Errorf("store get rows error: %w", err)
	}

	return order, nil
}

func (s *Store) ProcessOrder(orderNumber string, accrual float64) error {
	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE orders SET status = 'PROCESSED', accrual = $1 WHERE order_number = $2", accrual, orderNumber)
	if err != nil {
		return fmt.Errorf("db update error: %w", err)
	}

	order, err := s.GetOrdersUser(orderNumber)
	if err != nil {
		return fmt.Errorf("store user balance error: %w", err)
	}
	username := order.Username

	balance, err := s.GetUserBalance(username)
	if err != nil {
		return fmt.Errorf("store user balance error: %w", err)
	}

	_, err = db.Exec("UPDATE users SET balance = $1 WHERE username = $2", balance.Balance+accrual, username)
	if err != nil {
		return fmt.Errorf("store query error: %w", err)
	}

	return nil
}

func (s *Store) OrdersProcessing() ([]Order, error) {
	var result []Order

	db, err := sqlx.Connect("postgres", s.DB)
	if err != nil {
		return result, fmt.Errorf("store connect error: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE orders SET status = 'PROCESSING' WHERE status = 'NEW'")
	if err != nil {
		return result, fmt.Errorf("db update error: %w", err)
	}

	order := Order{}
	rows, err := db.Queryx("SELECT * FROM orders WHERE status = 'PROCESSING'")
	if err != nil {
		return result, fmt.Errorf("store query rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&order)
		if err != nil {
			return result, fmt.Errorf("store scan error: %w", err)
		}
		result = append(result, order)
	}

	err = rows.Err()
	if err != nil {
		return result, fmt.Errorf("store get rows error: %w", err)
	}

	return result, nil
}
