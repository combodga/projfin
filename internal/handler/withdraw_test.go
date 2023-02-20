package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/service"
	service_mocks "github.com/combodga/projfin/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
)

func TestHandler_PostBalanceWithdraw(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockWithdraw, ctx context.Context, username, orderNumber string, sum float64)

	tests := []struct {
		name                 string
		inputBody            string
		username             string
		orderNumber          string
		sum                  float64
		mockBehavior         mockBehavior
		expectedOrderStatus  projfin.OrderStatus
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			username:    "test",
			orderNumber: "1234567",
			sum:         1,
			mockBehavior: func(r *service_mocks.MockWithdraw, ctx context.Context, username, orderNumber string, sum float64) {
				r.EXPECT().Withdraw(context.Background(), username, orderNumber, sum).Return(projfin.OrderStatusOK, nil)
			},
			expectedOrderStatus:  projfin.OrderStatusOK,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "status: ok",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := service_mocks.NewMockWithdraw(c)
			test.mockBehavior(store, context.Background(), test.username, test.orderNumber, test.sum)

			services := &service.Service{Withdraw: store}
			handler := Handler{services}

			e := echo.New()
			e.POST("/api/user/balance/withdraw", handler.PostRegister)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/balance/withdraw",
				bytes.NewBufferString(test.inputBody))

			e.ServeHTTP(w, req)

			withdraw, err := services.Withdraw.Withdraw(context.Background(), test.username, test.orderNumber, test.sum)

			assert.Equal(t, withdraw, test.expectedOrderStatus)
			assert.Equal(t, err, nil)
		})
	}
}
