package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
)

type WithdrawService struct {
	sw store.Withdraw
}

func NewWithdrawService(sw store.Withdraw) *WithdrawService {
	return &WithdrawService{sw: sw}
}

func (s *WithdrawService) ListWithdrawals(username string) ([]projfin.WithdrawalsListItem, error) {
	var withdrawals []projfin.WithdrawalsListItem

	w, err := s.sw.ListWithdrawals(username)
	if err != nil {
		return withdrawals, fmt.Errorf("list withdrawals service error: %w", err)
	}

	for _, withdraw := range w {
		withdrawals = append(withdrawals, projfin.WithdrawalsListItem{OrderNum: withdraw.OrderNumber, Sum: withdraw.Sum, ProcessedAt: withdraw.ProcessedAt})
	}

	return withdrawals, nil
}

func (s *WithdrawService) Withdraw(ctx context.Context, username, orderNumber string, sum float64) (projfin.OrderStatus, error) {
	orderNumberInt, err := strconv.Atoi(orderNumber)
	if err != nil {
		return projfin.OrderStatusNotANumber, fmt.Errorf("withdraw service error: %w", err)
	}

	if !user.ValidateOrderNumber(orderNumberInt) {
		return projfin.OrderStatusNotValid, nil
	}

	withdraw, err := s.sw.Withdraw(ctx, username, orderNumber, sum)
	if err != nil {
		return projfin.OrderStatusError, fmt.Errorf("withdraw service error: %w", err)
	}

	if withdraw == http.StatusPaymentRequired {
		return projfin.OrderStatusPaymentRequired, nil
	}

	return projfin.OrderStatusOK, nil
}
