package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
)

type balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type withdraw struct {
	OrderNum string  `json:"order"`
	Sum      float64 `json:"sum"`
}

type withdrawalsList struct {
	OrderNum    string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (h *Handler) PostBalanceWithdraw(c echo.Context) error {
	if c.Get("username") == nil {
		return fmt.Errorf("get username error")
	}

	username := c.Get("username").(string)

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var wdraw withdraw
	err = json.Unmarshal(body, &wdraw)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	orderNum, err := strconv.Atoi(wdraw.OrderNum)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	if !user.ValidateOrderNumber(orderNum) {
		return c.String(http.StatusUnprocessableEntity, "status: unprocessable entity")
	}

	withdraw, err := h.Store.Withdraw.Withdraw(c.Request().Context(), username, wdraw.OrderNum, wdraw.Sum)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if withdraw == 402 {
		return c.String(http.StatusPaymentRequired, "status: payment required")
	}

	return c.String(http.StatusOK, "status: ok")
}

func (h *Handler) GetBalance(c echo.Context) error {
	if c.Get("username") == nil {
		return fmt.Errorf("get username error")
	}

	username := c.Get("username").(string)

	bal, err := h.Store.Order.GetUserBalance(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	return c.JSON(http.StatusOK, balance{bal.Balance, bal.Withdrawn})
}

func (h *Handler) GetWithdrawals(c echo.Context) error {
	if c.Get("username") == nil {
		return fmt.Errorf("get username error")
	}

	username := c.Get("username").(string)

	w, err := h.Store.Withdraw.ListWithdrawals(c.Request().Context(), username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if w == nil {
		return c.String(http.StatusNoContent, "status: no content")
	}

	var withdrawals []withdrawalsList
	for _, withdraw := range w {
		withdrawals = append(withdrawals, withdrawalsList{withdraw.OrderNumber, withdraw.Sum, withdraw.ProcessedAt})
	}

	return c.JSON(http.StatusOK, withdrawals)
}
