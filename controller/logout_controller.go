package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

func Logout(c echo.Context) error {
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	session.Clear(c)
	return c.Redirect(http.StatusFound, "/login")
}
