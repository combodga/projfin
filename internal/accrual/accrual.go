package accrual

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
