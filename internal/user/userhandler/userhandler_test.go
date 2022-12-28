package userhandler

import (
	"crypto/rand"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/combodga/projfin/internal/handler"
	"github.com/labstack/echo/v4"
)

func Test_PostRegister(t *testing.T) {
	h, err := handler.New("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "")
	if err != nil {
		t.Fatalf("error initializing: %v", err)
	}
	uh := New(h)

	e := echo.New()
	n, _ := rand.Int()
	s := "{\"login\":\"" + strconv.Itoa(n) + "\",\"password\":\"test\"}"
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", strings.NewReader(s))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	uh.PostRegister(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status %v; got %v", http.StatusOK, result.StatusCode)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("error parsing response: %v", err)
	}

	if string(body) != "status: ok" {
		t.Fatalf("expected answer to be %v; got %v", "status: ok", string(body))
	}
}

func Test_PostLogin(t *testing.T) {
	h, err := handler.New("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "")
	if err != nil {
		t.Fatalf("error initializing: %v", err)
	}
	uh := New(h)

	e := echo.New()
	s := "{\"login\":\"check\",\"password\":\"check\"}"
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", strings.NewReader(s))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	uh.PostRegister(c)

	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/login", strings.NewReader(s))

	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)
	uh.PostLogin(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status %v; got %v", http.StatusOK, result.StatusCode)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("error parsing response: %v", err)
	}

	if string(body) != "status: ok" {
		t.Errorf("expected answer to be %v; got %v", "status: ok", string(body))
	}

	s = "{\"login\":\"check\",\"password\":\"error\"}"
	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/login", strings.NewReader(s))

	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)
	uh.PostLogin(c)

	result = recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %v; got %v", http.StatusUnauthorized, result.StatusCode)
	}

	body, err = io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("error parsing response: %v", err)
	}

	if string(body) != "status: unathorized" {
		t.Errorf("expected answer to be %v; got %v", "status: unathorized", string(body))
	}
}
