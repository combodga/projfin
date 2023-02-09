package service

import (
	"context"
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

func (s *OrderService) MakeOrder(ctx context.Context, username, orderNumber string) error {
	err := s.so.MakeOrder(ctx, username, orderNumber)
	if err != nil {
		log.Printf("make order service error: %v", err)
	}
	return err
}

func (s *OrderService) ListOrders(username string) ([]projfin.OrderListItem, error) {
	orders, err := s.so.ListOrders(username)
	if err != nil {
		log.Printf("list orders service error: %v", err)
	}

	var ordersList []projfin.OrderListItem
	for _, o := range orders {
		ordersList = append(ordersList, projfin.OrderListItem{Number: o.OrderNumber, Status: o.Status, Accrual: o.Accrual, UploadedAt: o.UploadedAt})
	}

	return ordersList, err
}

func (s *OrderService) InvalidateOrder(orderNumber string) error {
	err := s.so.InvalidateOrder(orderNumber)
	if err != nil {
		log.Printf("invalidate order service error: %v", err)
	}
	return err
}

func (s *OrderService) GetOrdersUser(orderNumber string) (projfin.Order, error) {
	order, err := s.so.GetOrdersUser(orderNumber)
	if err != nil {
		log.Printf("get orders user service error: %v", err)
	}
	return order, err
}

func (s *OrderService) ProcessOrder(orderNumber string, accrual float64) error {
	err := s.so.ProcessOrder(orderNumber, accrual)
	if err != nil {
		log.Printf("process order service error: %v", err)
	}
	return err
}

func (s *OrderService) OrdersProcessing() ([]projfin.Order, error) {
	orders, err := s.so.OrdersProcessing()
	if err != nil {
		log.Printf("orders processing service error: %v", err)
	}
	return orders, err
}

func (s *OrderService) GetUserBalance(username string) (projfin.User, error) {
	user, err := s.so.GetUserBalance(username)
	if err != nil {
		log.Printf("get user balance service error: %v", err)
	}
	return user, err
}
