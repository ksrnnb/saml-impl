package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	return c.Render(http.StatusOK, "home.html", nil)
}
