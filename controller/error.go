package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

func errorRedirectToLogin(c echo.Context, msg string) error {
	session.SetFlash(c, "error", msg)
	return errorRedirect(c, "/login")
}

func errorRedirectToSAMLLogin(c echo.Context, msg string) error {
	session.SetFlash(c, "error", msg)
	return errorRedirect(c, "/login/saml")
}

func errorRedirect(c echo.Context, path string) error {
	return c.Redirect(http.StatusFound, path)
}
