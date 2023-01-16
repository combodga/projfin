package store

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Order interface {
	CheckOrder(username, orderNumber string) (int, error)
	MakeOrder(ctx context.Context, username, orderNumber string) error
	ListOrders(username string) ([]order, error)
	InvalidateOrder(orderNumber string) error
	GetOrdersUser(orderNumber string) (order, error)
	ProcessOrder(orderNumber string, accrual float64) error
	OrdersProcessing() ([]order, error)
	GetUserBalance(username string) (user, error)
}

type User interface {
	DoRegister(ctx context.Context, username, password string) error
	DoLogin(ctx context.Context, username, password string) error
}

type Withdraw interface {
	ListWithdrawals(ctx context.Context, username string) ([]withdraw, error)
	Withdraw(ctx context.Context, username, orderNumber string, sum float64) (int, error)
}

type Store struct {
	DB        *sqlx.DB
	ErrorDupe error
	Order
	User
	Withdraw
}

func New(database string) (*Store, error) {
	db, err := sqlx.Connect("postgres", database)
	if err != nil {
		return nil, fmt.Errorf("store connect error: %w", err)
	}

	ed := fmt.Errorf("duplicate key error")
	s := &Store{
		DB:        db,
		ErrorDupe: ed,
		Order:     NewOrderPG(db, ed),
		User:      NewUserPG(db, ed),
		Withdraw:  NewWithdrawPG(db),
	}

	sql := "BEGIN;CREATE TABLE IF NOT EXISTS users (username text primary key,password text,balance double precision,withdrawn double precision);"
	sql += "CREATE TABLE IF NOT EXISTS orders (order_number text primary key,username text REFERENCES users,status text,accrual double precision,uploaded_at timestamp with time zone);"
	sql += "CREATE TABLE IF NOT EXISTS withdrawals (order_number text,username text REFERENCES users,sum double precision,processed_at timestamp with time zone);COMMIT;"
	db.MustExec(sql)

	return s, nil
}

func (s *Store) Close() {
	s.DB.Close()
}
