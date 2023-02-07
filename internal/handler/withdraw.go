package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/combodga/projfin"
	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
)

func (h *Handler) PostBalanceWithdraw(c echo.Context) error {
	if c.Get("username") == nil {
		return fmt.Errorf("get username error")
	}

	username := c.Get("username").(string)

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var wdraw projfin.WithdrawShort
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

	projfin.Context = c.Request().Context()
	withdraw, err := h.services.Withdraw.Withdraw(username, wdraw.OrderNum, wdraw.Sum)
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

	bal, err := h.services.Order.GetUserBalance(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	return c.JSON(http.StatusOK, projfin.Balance{Current: bal.Balance, Withdrawn: bal.Withdrawn})
}

func (h *Handler) GetWithdrawals(c echo.Context) error {
	if c.Get("username") == nil {
		return fmt.Errorf("get username error")
	}

	username := c.Get("username").(string)

	w, err := h.services.Withdraw.ListWithdrawals(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if w == nil {
		return c.String(http.StatusNoContent, "status: no content")
	}

	var withdrawals []projfin.WithdrawalsList
	for _, withdraw := range w {
		withdrawals = append(withdrawals, projfin.WithdrawalsList{OrderNum: withdraw.OrderNumber, Sum: withdraw.Sum, ProcessedAt: withdraw.ProcessedAt})
	}

	return c.JSON(http.StatusOK, withdrawals)
}
