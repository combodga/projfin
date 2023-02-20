package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/combodga/projfin/internal/service"
	service_mocks "github.com/combodga/projfin/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
)

func TestHandler_PostOrders(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockOrder, ctx context.Context, username, orderNumber string)

	tests := []struct {
		name                 string
		username             string
		orderNumber          string
		mockBehavior         mockBehavior
		expectedOrderError   error
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			orderNumber: "1234567",
			mockBehavior: func(r *service_mocks.MockOrder, ctx context.Context, username, orderNumber string) {
				r.EXPECT().MakeOrder(context.Background(), username, orderNumber).Return(nil)
			},
			expectedOrderError:   nil,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "status: ok",
		},
		{
			name:        "Error",
			orderNumber: "",
			mockBehavior: func(r *service_mocks.MockOrder, ctx context.Context, username, orderNumber string) {
				r.EXPECT().MakeOrder(context.Background(), username, orderNumber).Return(fmt.Errorf("some error"))
			},
			expectedOrderError:   fmt.Errorf("some error"),
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "status: ok",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := service_mocks.NewMockOrder(c)
			test.mockBehavior(store, context.Background(), test.username, test.orderNumber)

			services := &service.Service{Order: store}
			handler := Handler{services}

			e := echo.New()
			e.POST("/api/user/orders", handler.PostOrders)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/orders",
				bytes.NewBufferString(test.orderNumber))

			e.ServeHTTP(w, req)

			status := services.Order.MakeOrder(context.Background(), test.username, test.orderNumber)

			assert.Equal(t, status, test.expectedOrderError)
		})
	}
}
