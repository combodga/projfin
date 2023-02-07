package service

import (
	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

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
