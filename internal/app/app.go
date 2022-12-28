package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/combodga/projfin/internal/accrual"
	"github.com/combodga/projfin/internal/handler"
	"github.com/combodga/projfin/internal/order/orderHandler"
	"github.com/combodga/projfin/internal/user/userHandler"
	"github.com/combodga/projfin/internal/withdraw/withdrawHandler"

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

	uh := userHandler.New(h)
	e.POST("/api/user/register", uh.PostRegister)
	e.POST("/api/user/login", uh.PostLogin)

	oh := orderHandler.New(h)
	e.POST("/api/user/orders", oh.PostOrders, handler.Auth)
	e.GET("/api/user/orders", oh.GetOrders, handler.Auth)

	wh := withdrawHandler.New(h)
	e.POST("/api/user/balance/withdraw", wh.PostBalanceWithdraw, handler.Auth)
	e.GET("/api/user/balance", wh.GetBalance, handler.Auth)
	e.GET("/api/user/withdrawals", wh.GetWithdrawals, handler.Auth)

	go accrual.FetchAccruals(h)

	go func() {
		if err := e.Start(run); err != nil && err != http.ErrServerClosed {
			h.Store.Close()
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		h.Store.Close()
		e.Logger.Fatal(err)
	}

	return nil
}
