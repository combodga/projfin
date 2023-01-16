package handler

import (
	"fmt"
	"net/http"

	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Store *store.Store
	DB    string
	Accr  string
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

// Middleware

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, err := user.Get(c)
		if err != nil {
			return c.String(http.StatusUnauthorized, "status: unauthorized")
		}
		c.Set("username", username)
		return next(c)
	}
}
