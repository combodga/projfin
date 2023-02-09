package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/combodga/projfin"
	"github.com/labstack/echo/v4"
)

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
	orderStatus := h.services.Order.CheckOrder(username, order)
	switch orderStatus {
	case projfin.OrderStatusNotANumber:
		return c.String(http.StatusBadRequest, "status: bad request")
	case projfin.OrderStatusNotValid:
		return c.String(http.StatusUnprocessableEntity, "status: unprocessable entity")
	case projfin.OrderStatusError:
		return c.String(http.StatusInternalServerError, "status: internal server error")
	case projfin.OrderStatusOccupied:
		return c.String(http.StatusConflict, "status: conflict")
	case projfin.OrderStatusExists:
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

	return c.JSON(http.StatusOK, orders)
}
