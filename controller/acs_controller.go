package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/service"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

// HTTP POST Binding
func ConsumeAssertion(c echo.Context) error {
	ss, err := service.NewSamlService(c.Param("company_id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	r := c.Request()
	r.ParseForm()

	// TODO: handle SP-initiated
	possibleRequestIDs := []string{}
	if ss.ServiceProvider.AllowIDPInitiated {
		possibleRequestIDs = append(possibleRequestIDs, "")
	}

	assertion, err := ss.ServiceProvider.ParseResponse(r, possibleRequestIDs)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// SAML 認証成功後
	email := assertion.Subject.NameID.Value
	u, err := model.FindUserByEmail(email)
	if err != nil {
		return c.String(http.StatusInternalServerError, "ユーザーがみつかりませんでした")
	}

	if err := session.Clear(c); err != nil {
		session.Set(c, "error", "予期しないエラーが発生しました")
		return c.Redirect(http.StatusFound, "/login")
	}

	session.Set(c, "userId", u.ID)
	session.Set(c, "success", "SAML 認証に成功しました")
	session.Activate(u.ID)

	redirect := "/"
	rs := c.FormValue("RelayState")
	if rs != "" {
		// TODO: handle RelayState
	}
	return c.Redirect(http.StatusFound, redirect)
}
