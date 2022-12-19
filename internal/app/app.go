package app

import (
	"fmt"

	"github.com/combodga/projfin/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Go(run, database, accr string) error {
	h, err := handler.New(database, accr)
	if err != nil {
		return fmt.Errorf("error handler init: %w", err)
	}

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	e.POST("/api/user/register", h.PostRegister)
	e.POST("/api/user/login", h.PostLogin)
	e.POST("/api/user/orders", h.PostOrders)
	e.POST("/api/user/balance/withdraw", h.PostBalanceWithdraw)

	e.GET("/api/user/orders", h.GetOrders)
	e.GET("/api/user/balance", h.GetBalance)
	e.GET("/api/user/withdrawals", h.GetWithdrawals)

	go h.FetchAccruals()
	e.Logger.Fatal(e.Start(run))

	return nil
}
