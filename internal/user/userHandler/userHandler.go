package userHandler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/combodga/projfin/internal/handler"
	"github.com/combodga/projfin/internal/store"
	"github.com/combodga/projfin/internal/user"
	"github.com/combodga/projfin/internal/user/userStore"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Store *store.Store
	DB    string
	Accr  string
}

type credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func New(h *handler.Handler) *UserHandler {
	return &UserHandler{
		Store: h.Store,
		DB:    h.DB,
		Accr:  h.Accr,
	}
}

func (uh *UserHandler) PostRegister(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var cred credentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	err = userStore.DoRegister(uh.Store, c.Request().Context(), cred.Username, user.PasswordHasher(cred.Password))
	if errors.Is(err, uh.Store.ErrorDupe) {
		return c.String(http.StatusConflict, "status: conflict")
	} else if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	user.Set(c, cred.Username)
	return c.String(http.StatusOK, "status: ok")
}

func (uh *UserHandler) PostLogin(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var cred credentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	hash := user.PasswordHasher(cred.Password)
	err = userStore.DoLogin(uh.Store, c.Request().Context(), cred.Username, hash)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unathorized")
	}

	user.Set(c, cred.Username)
	return c.String(http.StatusOK, "status: ok")
}
