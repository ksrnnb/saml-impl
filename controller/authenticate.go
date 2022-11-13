package controller

import (
	"net/http"

	"github.com/ksrnnb/saml/model"
	"github.com/ksrnnb/saml/session"
	"github.com/labstack/echo/v4"
)

func authenticate(c echo.Context) (*model.User, error) {
	uid, err := session.Get(c, "userId")
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, err.Error())
	}
	u := model.FindUser(uid)
	if u == nil {
		return nil, c.Redirect(http.StatusFound, "/login")
	}
	return u, nil
}

func notAuthenticate(c echo.Context) error {
	uid, err := session.Get(c, "userId")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if uid != "" {
		return c.Redirect(http.StatusFound, "/")
	}
	return nil
}
