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
	"github.com/combodga/projfin/internal/service"
	"github.com/combodga/projfin/internal/store"
)

func Go(run, database, accr string) error {
	db, err := store.NewPGDB(database)
	if err != nil {
		return fmt.Errorf("error initializing db: %w", err)
	}

	stores := store.NewStore(db)
	services := service.NewService(stores)
	handlers := handler.NewHandler(services)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go accrual.FetchAccruals(ctx, accr, stores)

	e := handlers.InitRoutes()

	go func() {
		if err := e.Start(run); err != nil && err != http.ErrServerClosed {
			store.PGClose(db)
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()
	if err := e.Shutdown(ctx2); err != nil {
		store.PGClose(db)
		e.Logger.Fatal(err)
	}

	return nil
}
