package service

import (
	"github.com/combodga/projfin/internal/store"
)

type UserService struct {
	su store.User
}

func NewUserService(su store.User) *UserService {
	return &UserService{su: su}
}

func (s *UserService) DoRegister(username, password string) error {
	return s.su.DoRegister(username, password)
}

func (s *UserService) DoLogin(username, password string) error {
	return s.su.DoLogin(username, password)
}
