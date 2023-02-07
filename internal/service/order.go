package service

import (
	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
)

type OrderService struct {
	so store.Order
}

func NewOrderService(so store.Order) *OrderService {
	return &OrderService{so: so}
}

func (s *OrderService) CheckOrder(username, orderNumber string) (int, error) {
	return s.so.CheckOrder(username, orderNumber)
}

func (s *OrderService) MakeOrder(username, orderNumber string) error {
	return s.so.MakeOrder(username, orderNumber)
}

func (s *OrderService) ListOrders(username string) ([]projfin.Order, error) {
	return s.so.ListOrders(username)
}

func (s *OrderService) InvalidateOrder(orderNumber string) error {
	return s.so.InvalidateOrder(orderNumber)
}

func (s *OrderService) GetOrdersUser(orderNumber string) (projfin.Order, error) {
	return s.so.GetOrdersUser(orderNumber)
}

func (s *OrderService) ProcessOrder(orderNumber string, accrual float64) error {
	return s.so.ProcessOrder(orderNumber, accrual)
}

func (s *OrderService) OrdersProcessing() ([]projfin.Order, error) {
	return s.so.OrdersProcessing()
}

func (s *OrderService) GetUserBalance(username string) (projfin.User, error) {
	return s.so.GetUserBalance(username)
}
