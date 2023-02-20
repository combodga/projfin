package service

import (
	"context"
	"fmt"

	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
)

type UserService struct {
	su store.User
}

func NewUserService(su store.User) *UserService {
	return &UserService{su: su}
}

func (s *UserService) DoRegister(ctx context.Context, username, password string) error {
	hash := user.PasswordHasher(password)
	err := s.su.DoRegister(ctx, username, hash)
	if err != nil {
		return fmt.Errorf("do register service error: %w", err)
	}

	return nil
}

func (s *UserService) DoLogin(ctx context.Context, username, password string) error {
	hash := user.PasswordHasher(password)
	err := s.su.DoLogin(ctx, username, hash)
	if err != nil {
		return fmt.Errorf("do login service error: %w", err)
	}

	return nil
}
