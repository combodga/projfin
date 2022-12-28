package orderhandler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/combodga/projfin/internal/handler"
	"github.com/combodga/projfin/internal/order/orderstore"
	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	Store *store.Store
	DB    string
	Accr  string
}

type order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

func New(h *handler.Handler) *OrderHandler {
	return &OrderHandler{
		Store: h.Store,
		DB:    h.DB,
		Accr:  h.Accr,
	}
}

func (oh *OrderHandler) PostOrders(c echo.Context) error {
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

	code, err := orderstore.CheckOrder(oh.Store, c.Request().Context(), username, order)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}
	if code == 2 {
		return c.String(http.StatusConflict, "status: conflict")
	} else if code == 1 {
		return c.String(http.StatusOK, "status: ok")
	}

	err = orderstore.MakeOrder(oh.Store, c.Request().Context(), username, order)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	return c.String(http.StatusAccepted, "status: accepted")
}

func (oh *OrderHandler) GetOrders(c echo.Context) error {
	username := c.Get("username").(string)

	orders, err := orderstore.ListOrders(oh.Store, username)
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
