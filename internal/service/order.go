package service

import (
	"context"
	"fmt"
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

func (s *OrderService) CheckOrder(username, orderNumber string) (projfin.OrderStatus, error) {
	orderNumberInt, err := strconv.Atoi(orderNumber)
	if err != nil {
		return projfin.OrderStatusNotANumber, fmt.Errorf("order number int error: %w", err)
	}

	if !user.ValidateOrderNumber(orderNumberInt) {
		return projfin.OrderStatusNotValid, nil
	}

	orderStatus, err := s.so.CheckOrder(username, orderNumber)
	if err != nil {
		return projfin.OrderStatusError, fmt.Errorf("check order error: %w", err)
	}

	return orderStatus, nil
}

func (s *OrderService) MakeOrder(ctx context.Context, username, orderNumber string) error {
	err := s.so.MakeOrder(ctx, username, orderNumber)
	if err != nil {
		return fmt.Errorf("make order service error: %w", err)
	}
	return nil
}

func (s *OrderService) ListOrders(username string) ([]projfin.OrderListItem, error) {
	var ordersList []projfin.OrderListItem

	orders, err := s.so.ListOrders(username)
	if err != nil {
		return ordersList, fmt.Errorf("list orders service error: %w", err)
	}

	for _, o := range orders {
		ordersList = append(ordersList, projfin.OrderListItem{Number: o.OrderNumber, Status: o.Status, Accrual: o.Accrual, UploadedAt: o.UploadedAt})
	}

	return ordersList, nil
}

func (s *OrderService) InvalidateOrder(orderNumber string) error {
	err := s.so.InvalidateOrder(orderNumber)
	if err != nil {
		return fmt.Errorf("invalidate order service error: %w", err)
	}
	return nil
}

func (s *OrderService) GetOrdersUser(orderNumber string) (projfin.Order, error) {
	order, err := s.so.GetOrdersUser(orderNumber)
	if err != nil {
		return order, fmt.Errorf("get orders user service error: %w", err)
	}
	return order, nil
}

func (s *OrderService) ProcessOrder(orderNumber string, accrual float64) error {
	err := s.so.ProcessOrder(orderNumber, accrual)
	if err != nil {
		return fmt.Errorf("process order service error: %w", err)
	}
	return nil
}

func (s *OrderService) OrdersProcessing() ([]projfin.Order, error) {
	orders, err := s.so.OrdersProcessing()
	if err != nil {
		return orders, fmt.Errorf("orders processing service error: %w", err)
	}
	return orders, err
}

func (s *OrderService) GetUserBalance(username string) (projfin.User, error) {
	user, err := s.so.GetUserBalance(username)
	if err != nil {
		return user, fmt.Errorf("get user balance service error: %w", err)
	}
	return user, nil
}
