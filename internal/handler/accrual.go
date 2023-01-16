package handler

import (
	"fmt"
	"time"

	"github.com/combodga/projfin/internal/accrual"
)

func (h *Handler) FetchAccruals() error {
	for {
		h.getAccruals()
		time.Sleep(300 * time.Millisecond)
	}
}

func (h *Handler) getAccruals() error {
	orders, err := h.Store.Order.OrdersProcessing()
	if err != nil {
		return fmt.Errorf("update accrual error: %w", err)
	}

	for _, order := range orders {
		status, accrual, err := accrual.Calculate(h.Accr, order.OrderNumber)
		if err != nil {
			return fmt.Errorf("update accrual order error: %w", err)
		}

		if status == "INVALID" {
			err = h.Store.Order.InvalidateOrder(order.OrderNumber)
			if err != nil {
				return fmt.Errorf("set order invalid error: %w", err)
			}
		} else if status == "PROCESSED" {
			err = h.Store.Order.ProcessOrder(order.OrderNumber, accrual)
			if err != nil {
				return fmt.Errorf("set order processed error: %w", err)
			}
		}
	}

	return nil
}
