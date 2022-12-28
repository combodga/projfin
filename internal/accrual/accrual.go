package accrual

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/combodga/projfin/internal/handler"
	"github.com/combodga/projfin/internal/order/orderStore"
)

type Accrual struct {
	OrderNum string  `json:"order"`
	Status   string  `json:"status"`
	Accrual  float64 `json:"accrual"`
}

func calculate(accr, orderNumber string) (string, float64, error) {
	resp, err := http.Get(accr + "/api/orders/" + orderNumber)
	if err != nil {
		return "", 0, fmt.Errorf("error getting accrual: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return "", 0, fmt.Errorf("error getting retry timeout: %w", err)
		}
		time.Sleep(time.Duration(retryAfter) * time.Second)
		return "", 0, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error getting accrual: %w", err)
	}

	var a Accrual
	err = json.Unmarshal(body, &a)
	if err != nil {
		return "", 0, fmt.Errorf("error getting accrual: %w", err)
	}

	return a.Status, a.Accrual, nil
}

func FetchAccruals(h *handler.Handler) error {
	for {
		getAccruals(h)
		time.Sleep(300 * time.Millisecond)
	}
}

func getAccruals(h *handler.Handler) error {
	orders, err := orderStore.OrdersProcessing(h.Store)
	if err != nil {
		return fmt.Errorf("update accrual error: %w", err)
	}

	for _, order := range orders {
		status, accrual, err := calculate(h.Accr, order.OrderNumber)
		if err != nil {
			return fmt.Errorf("update accrual order error: %w", err)
		}

		if status == "INVALID" {
			err = orderStore.InvalidateOrder(h.Store, order.OrderNumber)
			if err != nil {
				return fmt.Errorf("set order invalid error: %w", err)
			}
		} else if status == "PROCESSED" {
			err = orderStore.ProcessOrder(h.Store, order.OrderNumber, accrual)
			if err != nil {
				return fmt.Errorf("set order processed error: %w", err)
			}
		}
	}

	return nil
}
