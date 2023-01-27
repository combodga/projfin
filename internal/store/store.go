package store

import (
	"fmt"

	"github.com/combodga/projfin"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go

type Order interface {
	CheckOrder(username, orderNumber string) (int, error)
	MakeOrder(username, orderNumber string) error
	ListOrders(username string) ([]projfin.Order, error)
	InvalidateOrder(orderNumber string) error
	GetOrdersUser(orderNumber string) (projfin.Order, error)
	ProcessOrder(orderNumber string, accrual float64) error
	OrdersProcessing() ([]projfin.Order, error)
	GetUserBalance(username string) (projfin.User, error)
}

type User interface {
	DoRegister(username, password string) error
	DoLogin(username, password string) error
}

type Withdraw interface {
	ListWithdrawals(username string) ([]projfin.Withdraw, error)
	Withdraw(username, orderNumber string, sum float64) (int, error)
}

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
