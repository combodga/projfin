package service

import (
	"context"
	"fmt"
	"testing"

	store_mocks "github.com/combodga/projfin/internal/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_UserDoRegister(t *testing.T) {
	type mockBehavior func(r *store_mocks.MockUser, ctx context.Context, username, password string)

	tests := []struct {
		name         string
		username     string
		password     string
		mockBehavior mockBehavior
		awaitError   bool
	}{
		{
			name:     "Ok",
			username: "test",
			password: "test",
			mockBehavior: func(r *store_mocks.MockUser, ctx context.Context, username, password string) {
				r.EXPECT().DoRegister(context.Background(), username, password).Return(nil)
			},
		},
		{
			name:     "Error",
			username: "test",
			password: "test",
			mockBehavior: func(r *store_mocks.MockUser, ctx context.Context, username, password string) {
				r.EXPECT().DoRegister(context.Background(), username, password).Return(fmt.Errorf("some error"))
			},
			awaitError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			services := store_mocks.NewMockUser(c)
			test.mockBehavior(services, context.Background(), test.username, test.password)

			err := services.DoRegister(context.Background(), test.username, test.password)

			if test.awaitError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_UserDoLogin(t *testing.T) {
	type mockBehavior func(r *store_mocks.MockUser, ctx context.Context, username, password string)

	tests := []struct {
		name         string
		username     string
		password     string
		mockBehavior mockBehavior
		awaitError   bool
	}{
		{
			name:     "Ok",
			username: "test",
			password: "test",
			mockBehavior: func(r *store_mocks.MockUser, ctx context.Context, username, password string) {
				r.EXPECT().DoLogin(context.Background(), username, password).Return(nil)
			},
		},
		{
			name:     "Error",
			username: "test",
			password: "",
			mockBehavior: func(r *store_mocks.MockUser, ctx context.Context, username, password string) {
				r.EXPECT().DoLogin(context.Background(), username, password).Return(fmt.Errorf("some error"))
			},
			awaitError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			services := store_mocks.NewMockUser(c)
			test.mockBehavior(services, context.Background(), test.username, test.password)

			err := services.DoLogin(context.Background(), test.username, test.password)

			if test.awaitError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
