package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/store"
)

func FetchAccruals(ctx context.Context, accr string, stores *store.Store) error {
	for {
		getAccruals(accr, stores)
		// time.Sleep(300 * time.Millisecond)
		select {
		case <-time.After(300 * time.Millisecond):
			// pass
		case <-ctx.Done():
			return nil
		}
	}
}

func getAccruals(accr string, stores *store.Store) error {
	orders, err := stores.Order.OrdersProcessing()
	if err != nil {
		return fmt.Errorf("update accrual error: %w", err)
	}

	for _, order := range orders {
		status, accrual, err := calculate(accr, order.OrderNumber)
		if err != nil {
			return fmt.Errorf("update accrual order error: %w", err)
		}

		if status == "INVALID" {
			err = stores.Order.InvalidateOrder(order.OrderNumber)
			if err != nil {
				return fmt.Errorf("set order invalid error: %w", err)
			}
		} else if status == "PROCESSED" {
			err = stores.Order.ProcessOrder(order.OrderNumber, accrual)
			if err != nil {
				return fmt.Errorf("set order processed error: %w", err)
			}
		}
	}

	return nil
}

func calculate(accr string, orderNumber string) (string, float64, error) {
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

	var a projfin.Accrual
	err = json.Unmarshal(body, &a)
	if err != nil {
		return "", 0, fmt.Errorf("error getting accrual: %w", err)
	}

	return a.Status, a.Accrual, nil
}
