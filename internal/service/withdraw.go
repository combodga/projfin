package service

import (
	"context"

	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
)

type WithdrawService struct {
	sw store.Withdraw
}

func NewWithdrawService(sw store.Withdraw) *WithdrawService {
	return &WithdrawService{sw: sw}
}

func (s *WithdrawService) ListWithdrawals(ctx context.Context, username string) ([]projfin.Withdraw, error) {
	return s.sw.ListWithdrawals(ctx, username)
}

func (s *WithdrawService) Withdraw(ctx context.Context, username, orderNumber string, sum float64) (int, error) {
	return s.sw.Withdraw(ctx, username, orderNumber, sum)
}
