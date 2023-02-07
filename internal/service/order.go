package service

import (
	"log"
	"strconv"

	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
)

type OrderService struct {
	so store.Order
}

func NewOrderService(so store.Order) *OrderService {
	return &OrderService{so: so}
}

func (s *OrderService) CheckOrder(username, orderNumber string) projfin.OrderStatus {
	orderNumberInt, err := strconv.Atoi(orderNumber)
	if err != nil {
		return projfin.OrderStatusNotANumber
	}

	if !user.ValidateOrderNumber(orderNumberInt) {
		return projfin.OrderStatusNotValid
	}

	orderStatus, err := s.so.CheckOrder(username, orderNumber)
	if err != nil {
		log.Printf("check order error: %v", err)
		return projfin.OrderStatusError
	}

	return orderStatus
}

func (s *OrderService) MakeOrder(username, orderNumber string) error {
	return s.so.MakeOrder(username, orderNumber)
}

func (s *OrderService) ListOrders(username string) ([]projfin.OrderListItem, error) {
	orders, err := s.so.ListOrders(username)

	var ordersList []projfin.OrderListItem
	for _, o := range orders {
		ordersList = append(ordersList, projfin.OrderListItem{Number: o.OrderNumber, Status: o.Status, Accrual: o.Accrual, UploadedAt: o.UploadedAt})
	}

	return ordersList, err
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
