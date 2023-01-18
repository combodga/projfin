package handler

import (
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func Test_PostRegister(t *testing.T) {
	h, err := New("host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", "")
	if err != nil {
		t.Fatalf("error initializing: %v", err)
	}

	e := echo.New()
	rand.Seed(time.Now().UnixNano())

	s := "{\"login\":\"" + strconv.Itoa(rand.Intn(1e10)) + "\",\"password\":\"test\"}"
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", strings.NewReader(s))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	h.PostRegister(c)

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

	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/login", strings.NewReader(s))

	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)
	h.PostLogin(c)

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
	h.PostLogin(c)

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
