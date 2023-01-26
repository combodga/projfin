package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
)

type order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

func (h *Handler) PostOrders(c echo.Context) error {
	if c.Get("username") == nil {
		return fmt.Errorf("get username error")
	}

	username := c.Get("username").(string)

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

	code, err := h.services.Order.CheckOrder(username, order)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if code == 2 {
		return c.String(http.StatusConflict, "status: conflict")
	} else if code == 1 {
		return c.String(http.StatusOK, "status: ok")
	}

	err = h.services.Order.MakeOrder(c.Request().Context(), username, order)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	return c.String(http.StatusAccepted, "status: accepted")
}

func (h *Handler) GetOrders(c echo.Context) error {
	if c.Get("username") == nil {
		return fmt.Errorf("get username error")
	}

	username := c.Get("username").(string)

	orders, err := h.services.Order.ListOrders(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if orders == nil {
		return c.String(http.StatusNoContent, "status: no content")
	}

	var ordersList []order
	for _, o := range orders {
		ordersList = append(ordersList, order{o.OrderNumber, o.Status, o.Accrual, o.UploadedAt})
	}

	return c.JSON(http.StatusOK, ordersList)
}
