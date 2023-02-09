package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

	ctxAccruals := context.Background()
	ctxAccruals, cancelAccruals := context.WithCancel(ctxAccruals)
	defer cancelAccruals()

	var wg sync.WaitGroup
	wg.Add(1)

	go accrual.FetchAccruals(&wg, ctxAccruals, accr, stores)

	e := handlers.InitRoutes()

	go func() {
		if err := e.Start(run); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	wg.Wait()

	ctxServer, cancelServer := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelServer()

	if err := e.Shutdown(ctxServer); err != nil {
		log.Printf("server shutdown error: %v", err)
		cancelServer()
	}

	store.PGClose(db)

	return nil
}
