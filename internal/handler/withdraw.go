package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/combodga/projfin"
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

	projfin.Context = c.Request().Context()
	withdraw := h.services.Withdraw.Withdraw(username, wdraw.OrderNum, wdraw.Sum)
	switch withdraw {
	case projfin.OrderStatusNotANumber:
		return c.String(http.StatusBadRequest, "status: bad request")
	case projfin.OrderStatusNotValid:
		return c.String(http.StatusUnprocessableEntity, "status: unprocessable entity")
	case projfin.OrderStatusError:
		return c.String(http.StatusInternalServerError, "status: internal server error")
	case projfin.OrderStatusPaymentRequired:
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

	return c.JSON(http.StatusOK, w)
}
