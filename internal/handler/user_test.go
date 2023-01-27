package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/combodga/projfin/internal/service"
	service_mocks "github.com/combodga/projfin/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
)

func TestHandler_UserRegister(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockUser, username, password string)

	tests := []struct {
		name                 string
		inputBody            string
		username             string
		password             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "test", "password": "test"}`,
			username:  "test",
			password:  "dce472b679aa4d3893d3166dee95725a",
			mockBehavior: func(r *service_mocks.MockUser, username, password string) {
				r.EXPECT().DoRegister(username, password).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "status: ok",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := service_mocks.NewMockUser(c)
			test.mockBehavior(store, test.username, test.password)

			services := &service.Service{User: store}
			handler := Handler{services}

			e := echo.New()
			e.POST("/api/user/register", handler.PostRegister)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/register",
				bytes.NewBufferString(test.inputBody))

			e.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_UserLogin(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockUser, username, password string)

	tests := []struct {
		name                 string
		inputBody            string
		username             string
		password             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "test", "password": "test"}`,
			username:  "test",
			password:  "dce472b679aa4d3893d3166dee95725a",
			mockBehavior: func(r *service_mocks.MockUser, username, password string) {
				r.EXPECT().DoLogin(username, password).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "status: ok",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			store := service_mocks.NewMockUser(c)
			test.mockBehavior(store, test.username, test.password)

			services := &service.Service{User: store}
			handler := Handler{services}

			e := echo.New()
			e.POST("/api/user/login", handler.PostLogin)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/login",
				bytes.NewBufferString(test.inputBody))

			e.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
