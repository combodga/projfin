package handler

import (
	"net/http"

	"github.com/combodga/projfin/internal/service"
	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Middleware

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, err := user.GetAuthCookie(c)
		if err != nil {
			return c.String(http.StatusUnauthorized, "status: unauthorized")
		}
		c.Set("username", username)
		return next(c)
	}
}

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	g := e.Group("/api/user")

	g.POST("/register", h.PostRegister)
	g.POST("/login", h.PostLogin)

	g.POST("/orders", h.PostOrders, Auth)
	g.GET("/orders", h.GetOrders, Auth)

	g.POST("/balance/withdraw", h.PostBalanceWithdraw, Auth)
	g.GET("/balance", h.GetBalance, Auth)
	g.GET("/withdrawals", h.GetWithdrawals, Auth)

	return e
}
