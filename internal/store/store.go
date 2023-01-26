package store

import (
	"context"
	"fmt"

	"github.com/combodga/projfin"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go

type Order interface {
	CheckOrder(username, orderNumber string) (int, error)
	MakeOrder(ctx context.Context, username, orderNumber string) error
	ListOrders(username string) ([]projfin.Order, error)
	InvalidateOrder(orderNumber string) error
	GetOrdersUser(orderNumber string) (projfin.Order, error)
	ProcessOrder(orderNumber string, accrual float64) error
	OrdersProcessing() ([]projfin.Order, error)
	GetUserBalance(username string) (projfin.User, error)
}

type User interface {
	DoRegister(ctx context.Context, username, password string) error
	DoLogin(ctx context.Context, username, password string) error
}

type Withdraw interface {
	ListWithdrawals(ctx context.Context, username string) ([]projfin.Withdraw, error)
	Withdraw(ctx context.Context, username, orderNumber string, sum float64) (int, error)
}

// type Store struct {
// 	DB        *sqlx.DB
// 	ErrorDupe error
// 	Order
// 	User
// 	Withdraw
// }

type Store struct {
	Order
	User
	Withdraw
}

var (
	ErrorDupe = fmt.Errorf("duplicate key error")
)

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Order:    NewOrderPG(db),
		User:     NewUserPG(db),
		Withdraw: NewWithdrawPG(db),
	}
}
