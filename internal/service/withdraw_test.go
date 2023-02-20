package service

import (
	"fmt"
	"testing"

	"github.com/combodga/projfin"
	store_mocks "github.com/combodga/projfin/internal/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_WithdrawListWithdrawals(t *testing.T) {
	type mockBehavior func(r *store_mocks.MockWithdraw, username string)

	tests := []struct {
		name         string
		username     string
		mockBehavior mockBehavior
		awaitError   bool
	}{
		{
			name:     "Ok",
			username: "test",
			mockBehavior: func(r *store_mocks.MockWithdraw, username string) {
				r.EXPECT().ListWithdrawals(username).Return(make([]projfin.Withdraw, 0), nil)
			},
		},
		{
			name:     "Error",
			username: "",
			mockBehavior: func(r *store_mocks.MockWithdraw, username string) {
				r.EXPECT().ListWithdrawals(username).Return(make([]projfin.Withdraw, 0), fmt.Errorf("some error"))
			},
			awaitError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			services := store_mocks.NewMockWithdraw(c)
			test.mockBehavior(services, test.username)

			_, err := services.ListWithdrawals(test.username)

			if test.awaitError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
