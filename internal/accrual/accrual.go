package accrual

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Accrual struct {
	OrderNum string  `json:"order"`
	Status   string  `json:"status"`
	Accrual  float64 `json:"accrual"`
}

func Calculate(accr, orderNumber string) (string, float64, error) {
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
