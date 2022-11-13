package controller

import (
	"net/http"

	"github.com/ksrnnb/saml/session"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	msg, err := session.Get(c, "success")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "home.html", msg)
}
