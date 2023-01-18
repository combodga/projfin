package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func Test_GetBalance(t *testing.T) {
	h, err := New("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "")
	if err != nil {
		t.Fatalf("error initializing: %v", err)
	}

	e := echo.New()
	s := "{\"login\":\"check\",\"password\":\"check\"}"
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", strings.NewReader(s))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	h.PostRegister(c)

	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/orders", nil)

	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)
	c.Set("username", "check")
	h.GetBalance(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status %v; got %v", http.StatusOK, result.StatusCode)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("error parsing response: %v", err)
	}

	if string(body) == "0" {
		t.Errorf("expected answer to be %v; got %v", "0", string(body))
	}
}
