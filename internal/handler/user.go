package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/combodga/projfin/internal/user"
	"github.com/labstack/echo/v4"
)

type credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) PostRegister(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "status: internal server error")
	}

	var cred credentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	err = h.Store.User.DoRegister(c.Request().Context(), cred.Username, user.PasswordHasher(cred.Password))
	if errors.Is(err, h.Store.ErrorDupe) {
		return c.String(http.StatusConflict, "status: conflict")
	} else if err != nil {
		return c.String(http.StatusBadRequest, "status: bad request")
	}

	user.Set(c, cred.Username)
	return c.String(http.StatusOK, "status: ok")
}

func (h *Handler) PostLogin(c echo.Context) error {
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
	err = h.Store.User.DoLogin(c.Request().Context(), cred.Username, hash)
	if err != nil {
		return c.String(http.StatusUnauthorized, "status: unathorized")
	}

	user.Set(c, cred.Username)
	return c.String(http.StatusOK, "status: ok")
}
