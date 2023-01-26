package service

import (
	"context"

	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

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

type Service struct {
	Order
	User
	Withdraw
}

func NewService(stores *store.Store) *Service {
	return &Service{
		Order:    NewOrderService(stores.Order),
		User:     NewUserService(stores.User),
		Withdraw: NewWithdrawService(stores.Withdraw),
	}
}
