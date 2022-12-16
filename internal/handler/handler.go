package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/combodga/projfin/internal/accrual"
	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Store *store.Store
	DB    string
	Accr  string
}

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	OrderNum string  `json:"order"`
	Sum      float64 `json:"sum"`
}

type WithdrawalsList struct {
	OrderNum    string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func New(database, accr string) (*Handler, error) {
	s, err := store.New(database)
	if err != nil {
		err = fmt.Errorf("error store init: %w", err)
	}
	return &Handler{
		Store: s,
		DB:    database,
		Accr:  accr,
	}, err
}

// POST

func (h *Handler) PostRegister(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var cred Credentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	err = h.Store.DoRegister(cred.Username, user.PasswordHasher(cred.Password))
	if errors.Is(err, h.Store.ErrorDupe) {
		return c.String(http.StatusConflict, "status: conflict")
	} else if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	user.Set(c, cred.Username)
	return c.String(http.StatusOK, "status: ok")
}

func (h *Handler) PostLogin(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var cred Credentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	hash := user.PasswordHasher(cred.Password)
	err = h.Store.DoLogin(cred.Username, hash)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unathorized")
	}

	user.Set(c, cred.Username)
	return c.String(http.StatusOK, "status: ok")
}

func (h *Handler) PostOrders(c echo.Context) error {
	username, err := user.Get(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unauthorized")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	order := string(body)
	orderNum, err := strconv.Atoi(order)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	if !user.ValidateOrderNumber(orderNum) {
		return c.String(http.StatusUnprocessableEntity, "status: unprocessable entity")
	}

	code, err := h.Store.CheckOrder(username, order)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if code == 2 {
		return c.String(http.StatusConflict, "status: conflict")
	} else if code == 1 {
		return c.String(http.StatusOK, "status: ok")
	}

	err = h.Store.MakeOrder(username, order)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	return c.String(http.StatusAccepted, "status: accepted")
}

func (h *Handler) PostBalanceWithdraw(c echo.Context) error {
	username, err := user.Get(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unauthorized")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var wdraw Withdraw
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

	withdraw, err := h.Store.Withdraw(username, wdraw.OrderNum, wdraw.Sum)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if withdraw == 402 {
		return c.String(http.StatusPaymentRequired, "status: payment required")
	}

	return c.String(http.StatusOK, "status: ok")
}

// GET

func (h *Handler) GetOrders(c echo.Context) error {
	username, err := user.Get(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unauthorized")
	}

	orders, err := h.Store.ListOrders(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if orders == nil {
		return c.String(http.StatusNoContent, "status: no content")
	}

	var ordersList []Order
	for _, order := range orders {
		ordersList = append(ordersList, Order{order.OrderNumber, order.Status, order.Accrual, order.UploadedAt})
	}

	fmt.Println(username, ordersList)

	return c.JSON(http.StatusOK, ordersList)
}

func (h *Handler) GetBalance(c echo.Context) error {
	username, err := user.Get(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unauthorized")
	}

	balance, err := h.Store.GetUserBalance(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	return c.JSON(http.StatusOK, Balance{balance.Balance, balance.Withdrawn})
}

func (h *Handler) GetWithdrawals(c echo.Context) error {
	username, err := user.Get(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unauthorized")
	}

	withdrawals, err := h.Store.ListOrders(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if withdrawals == nil {
		return c.String(http.StatusNoContent, "status: no content")
	}

	var withdrawalsList []WithdrawalsList
	for _, withdraw := range withdrawals {
		withdrawalsList = append(withdrawalsList, WithdrawalsList{withdraw.OrderNumber, withdraw.Sum, withdraw.ProcessedAt})
	}

	return c.JSON(http.StatusOK, withdrawalsList)
}

func (h *Handler) GetAccruals() error {
	for {
		orders, err := h.Store.OrdersProcessing()
		if err != nil {
			return fmt.Errorf("update accrual error: %w", err)
		}

		for _, order := range orders {
			status, accrual, err := accrual.Calculate(h.Accr, order.OrderNumber)
			if err != nil {
				return fmt.Errorf("update accrual order error: %w", err)
			}
			if status == "INVALID" {
				err = h.Store.InvalidateOrder(order.OrderNumber)
				if err != nil {
					return fmt.Errorf("set order invalid error: %w", err)
				}
			} else if status == "PROCESSED" {
				err = h.Store.ProcessOrder(order.OrderNumber, accrual)
				if err != nil {
					return fmt.Errorf("set order processed error: %w", err)
				}
			}
		}
		time.Sleep(500)
	}
}
