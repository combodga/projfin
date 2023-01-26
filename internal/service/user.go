package service

import (
	"context"

	"github.com/combodga/projfin/internal/store"
)

type UserService struct {
	su store.User
}

func NewUserService(su store.User) *UserService {
	return &UserService{su: su}
}

func (s *UserService) DoRegister(ctx context.Context, username, password string) error {
	return s.su.DoRegister(ctx, username, password)
}

func (s *UserService) DoLogin(ctx context.Context, username, password string) error {
	return s.su.DoLogin(ctx, username, password)
}
