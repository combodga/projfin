package service

import (
	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
)

type WithdrawService struct {
	sw store.Withdraw
}

func NewWithdrawService(sw store.Withdraw) *WithdrawService {
	return &WithdrawService{sw: sw}
}

func (s *WithdrawService) ListWithdrawals(username string) ([]projfin.Withdraw, error) {
	return s.sw.ListWithdrawals(username)
}

func (s *WithdrawService) Withdraw(username, orderNumber string, sum float64) (int, error) {
	return s.sw.Withdraw(username, orderNumber, sum)
}
