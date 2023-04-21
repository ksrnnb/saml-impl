package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/service"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

const (
	StatusSeccess = "urn:oasis:names:tc:SAML:2.0:status:Success"
)

const (
	unexpectedError   = "予期しないエラーが発生しました"
	unexpectedMessage = "予期しないメッセージが送信されました"
	wrongIdentity     = "ユーザーIDまたはパスワードのいずれかが間違っています"
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
	errMsg, err := session.Get(c, "error")
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

type ShowCompanyLoginParams struct {
	ErrorMessage string
	Company      *model.Company
}

func ShowCompanyLogin(c echo.Context) error {
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}
	errMsg, err := session.Get(c, "error")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	company, err := model.FindCompany(c.Param("company_id"))
	if err != nil {
		return err
	}
	if company.IsZero() {
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Render(http.StatusOK, "company_login.html", ShowCompanyLoginParams{ErrorMessage: errMsg, Company: company})
}

func Login(c echo.Context) error {
	// NOTE: need to protect from csrf
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}

	// TODO: SAML 認証が有効なユーザーはパスワード認証できないようにする
	uid = c.FormValue("userId")
	pwd := c.FormValue("password")
	u, err := model.FindUser(uid)
	if err != nil {
		return errorRedirect(c, unexpectedError)
	}
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
	session.Activate(uid)
	return c.Redirect(http.StatusFound, "/")
}

// SAMLLogin starts to SP-initiated authentication.
func SAMLLogin(c echo.Context) error {
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}

	ss, err := service.NewSamlService(c.Param("company_id"))
	if err != nil {
		return err
	}
	u, err := ss.MakeAuthnRequestURL("")
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, u.String())
}

func errorRedirect(c echo.Context, msg string) error {
	session.Set(c, "error", msg)
	return c.Redirect(http.StatusFound, "/login")
}
