package controller

import (
	"net/http"

	"github.com/ksrnnb/saml/model"
	"github.com/ksrnnb/saml/session"
	"github.com/labstack/echo/v4"
)

const (
	StatusSeccess = "urn:oasis:names:tc:SAML:2.0:status:Success"
)

func ShowLogin(c echo.Context) error {
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}
	errMsg, err := session.Get(c, "error")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.Render(http.StatusOK, "login.html", errMsg)
}

func Login(c echo.Context) error {
	// NOTE: need to protect from csrf
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}

	uid = c.FormValue("userId")
	pwd := c.FormValue("password")
	u := model.FindUser(uid)

	if u == nil {
		session.Set(c, "error", "ユーザーIDまたはパスワードのいずれかが間違っています")
		return c.Redirect(http.StatusFound, "/login")
	}
	if err := u.ValidatePassword(pwd); err != nil {
		session.Set(c, "error", "ユーザーIDまたはパスワードのいずれかが間違っています")
		return c.Redirect(http.StatusFound, "/login")
	}
	// start to login
	if err := session.Clear(c); err != nil {
		session.Set(c, "error", "予期しないエラーが発生しました")
		return c.Redirect(http.StatusFound, "/login")
	}
	if err := session.Set(c, "userId", uid); err != nil {
		session.Set(c, "error", "予期しないエラーが発生しました")
		return c.Redirect(http.StatusFound, "/login")
	}

	return c.Redirect(http.StatusFound, "/")
}
