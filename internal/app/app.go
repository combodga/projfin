package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	e.POST("/api/user/orders", h.PostOrders, handler.Auth)
	e.GET("/api/user/orders", h.GetOrders, handler.Auth)

	e.POST("/api/user/balance/withdraw", h.PostBalanceWithdraw, handler.Auth)
	e.GET("/api/user/balance", h.GetBalance, handler.Auth)
	e.GET("/api/user/withdrawals", h.GetWithdrawals, handler.Auth)

	go h.FetchAccruals()

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
