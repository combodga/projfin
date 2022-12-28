package withdrawHandler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/combodga/projfin/internal/handler"
	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
	"github.com/combodga/projfin/internal/withdraw/withdrawStore"
	"github.com/labstack/echo/v4"
)

type WithdrawHandler struct {
	Store *store.Store
	DB    string
	Accr  string
}

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

func New(h *handler.Handler) *WithdrawHandler {
	return &WithdrawHandler{
		Store: h.Store,
		DB:    h.DB,
		Accr:  h.Accr,
	}
}

func (wh *WithdrawHandler) PostBalanceWithdraw(c echo.Context) error {
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

	withdraw, err := withdrawStore.Withdraw(wh.Store, c.Request().Context(), username, wdraw.OrderNum, wdraw.Sum)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if withdraw == 402 {
		return c.String(http.StatusPaymentRequired, "status: payment required")
	}

	return c.String(http.StatusOK, "status: ok")
}

func (wh *WithdrawHandler) GetBalance(c echo.Context) error {
	username := c.Get("username").(string)

	bal, err := withdrawStore.GetUserBalance(wh.Store, username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	return c.JSON(http.StatusOK, balance{bal.Balance, bal.Withdrawn})
}

func (wh *WithdrawHandler) GetWithdrawals(c echo.Context) error {
	username := c.Get("username").(string)

	w, err := withdrawStore.ListWithdrawals(wh.Store, c.Request().Context(), username)
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
