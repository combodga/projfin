package service

import (
	"fmt"
	"testing"

	"github.com/combodga/projfin"
	store_mocks "github.com/combodga/projfin/internal/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_OrderCheckOrder(t *testing.T) {
	type mockBehavior func(r *store_mocks.MockOrder, username, orderNumber string)

	tests := []struct {
		name         string
		username     string
		orderNumber  string
		mockBehavior mockBehavior
		awaitError   bool
	}{
		{
			name:        "Ok",
			username:    "test",
			orderNumber: "1234567",
			mockBehavior: func(r *store_mocks.MockOrder, username, orderNumber string) {
				r.EXPECT().CheckOrder(username, orderNumber).Return(projfin.OrderStatusOK, nil)
			},
		},
		{
			name:        "Error",
			username:    "test",
			orderNumber: "",
			mockBehavior: func(r *store_mocks.MockOrder, username, orderNumber string) {
				r.EXPECT().CheckOrder(username, orderNumber).Return(projfin.OrderStatusError, fmt.Errorf("some error"))
			},
			awaitError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			services := store_mocks.NewMockOrder(c)
			test.mockBehavior(services, test.username, test.orderNumber)

			status, err := services.CheckOrder(test.username, test.orderNumber)

			if test.awaitError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, status, projfin.OrderStatusOK)
				assert.NoError(t, err)
			}
		})
	}
}
