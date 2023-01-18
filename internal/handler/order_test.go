package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/combodga/projfin/internal/luhn"
	"github.com/labstack/echo/v4"
)

var rnd string

func Test_Init(t *testing.T) {
	// rand.Seed(time.Now().UnixNano())
	// rnd = strconv.Itoa(rand.Intn(1e10))
	rnd = luhn.GenerateLuhn(10)
}

func Test_PostOrders(t *testing.T) {
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

	s = rnd
	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/user/orders", strings.NewReader(s))

	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)
	c.Set("username", "check")
	h.PostOrders(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusAccepted {
		t.Errorf("expected status %v; got %v", http.StatusAccepted, result.StatusCode)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("error parsing response: %v", err)
	}

	if string(body) != "status: accepted" {
		t.Fatalf("expected answer to be %v; got %v", "status: accepted", string(body))
	}
}

func Test_GetOrders(t *testing.T) {
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

	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/orders", nil)

	recorder = httptest.NewRecorder()
	c = e.NewContext(request, recorder)
	h.GetOrders(c)

	result = recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status %v; got %v", http.StatusOK, result.StatusCode)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("error parsing response: %v", err)
	}

	if strings.Contains(string(body), "\""+rnd+"\"") {
		t.Errorf("expected answer must include %v; got %v", "\""+rnd+"\"", string(body))
	}
}
