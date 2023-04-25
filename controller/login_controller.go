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
	unexpectedError     = "予期しないエラーが発生しました"
	unexpectedMessage   = "予期しないメッセージが送信されました"
	wrongIdentity       = "ユーザーIDまたはパスワードのいずれかが間違っています"
	cannotPasswordLogin = "SAML 認証が有効なユーザーはパスワードでログインできません"
)

type ShowLoginArg struct {
	Message string
	Users   []*model.User
}

func ShowLogin(c echo.Context) error {
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}
	errMsg, err := session.GetFlash(c, "error")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	users, err := model.ListAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	arg := ShowLoginArg{
		Message: errMsg,
		Users:   users,
	}
	return c.Render(http.StatusOK, "login.html", arg)
}

func Login(c echo.Context) error {
	// NOTE: need to protect from csrf
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}

	uid = c.FormValue("userId")
	pwd := c.FormValue("password")
	u, err := model.FindUser(uid)
	if err != nil {
		return errorRedirectToLogin(c, unexpectedError)
	}
	if u == nil {
		return errorRedirectToLogin(c, wrongIdentity)
	}
	if !u.IsAdmin() {
		return errorRedirectToLogin(c, cannotPasswordLogin)
	}
	if err := u.ValidatePassword(pwd); err != nil {
		return errorRedirectToLogin(c, wrongIdentity)
	}

	// start to login
	if err := session.Clear(c); err != nil {
		return errorRedirectToLogin(c, unexpectedMessage)
	}
	sid, err := session.Start(c)
	if err != nil {
		return errorRedirectToLogin(c, unexpectedMessage)
	}
	session.Add(c, sid, "userId", uid)
	session.Activate(uid)
	return c.Redirect(http.StatusFound, "/")
}
