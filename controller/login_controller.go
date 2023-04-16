package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

const (
	StatusSeccess = "urn:oasis:names:tc:SAML:2.0:status:Success"
)

const (
	unexpectedMessage = "予期しないメッセージが送信されました"
	wrongIdentity     = "ユーザーIDまたはパスワードのいずれかが間違っています"
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
		return errorRedirect(c, wrongIdentity)
	}
	if err := u.ValidatePassword(pwd); err != nil {
		return errorRedirect(c, wrongIdentity)
	}

	// start to login
	if err := session.Clear(c); err != nil {
		return errorRedirect(c, unexpectedMessage)
	}
	if err := session.Set(c, "userId", uid); err != nil {
		return errorRedirect(c, unexpectedMessage)
	}

	return c.Redirect(http.StatusFound, "/")
}

func errorRedirect(c echo.Context, msg string) error {
	session.Set(c, "error", msg)
	return c.Redirect(http.StatusFound, "/login")
}
