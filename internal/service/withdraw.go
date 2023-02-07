package service

import (
	"log"
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
	w, err := s.sw.ListWithdrawals(username)

	var withdrawals []projfin.WithdrawalsListItem
	for _, withdraw := range w {
		withdrawals = append(withdrawals, projfin.WithdrawalsListItem{OrderNum: withdraw.OrderNumber, Sum: withdraw.Sum, ProcessedAt: withdraw.ProcessedAt})
	}

	return withdrawals, err
}

func (s *WithdrawService) Withdraw(username, orderNumber string, sum float64) projfin.OrderStatus {
	orderNumberInt, err := strconv.Atoi(orderNumber)
	if err != nil {
		return projfin.OrderStatusNotANumber
	}

	if !user.ValidateOrderNumber(orderNumberInt) {
		return projfin.OrderStatusNotValid
	}

	withdraw, err := s.sw.Withdraw(username, orderNumber, sum)
	if err != nil {
		log.Printf("withdraw service error: %v", err)
		return projfin.OrderStatusError
	}

	if withdraw == http.StatusPaymentRequired {
		return projfin.OrderStatusPaymentRequired
	}

	return projfin.OrderStatusOK
}
