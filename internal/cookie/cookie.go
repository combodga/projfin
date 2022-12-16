package cookie

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Write(c echo.Context, name, value string) {
	ckie := new(http.Cookie)
	ckie.Name = name
	ckie.Value = value
	c.SetCookie(ckie)
}

func Read(c echo.Context, name string) (string, error) {
	ckie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return ckie.Value, nil
}
