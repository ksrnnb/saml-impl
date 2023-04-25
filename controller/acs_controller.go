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
		return errorRedirectToLogin(c, err.Error())
	}

	r := c.Request()
	r.ParseForm()

	possibleRequestIDs := service.ListRequestIDs()
	if ss.ServiceProvider.AllowIDPInitiated {
		possibleRequestIDs = append(possibleRequestIDs, "")
	}

	assertion, err := ss.ParseReponse(r, possibleRequestIDs)
	if err != nil {
		return errorRedirectToLogin(c, err.Error())
	}

	// AllowIDPInitiated == true の場合は、 SP-initiated のリクエストが来ても
	// InResponseTo は検証しないようになっているので自前で検証する
	if err := ss.ValidateInResponseTo(c.FormValue("SAMLResponse")); err != nil {
		return errorRedirectToLogin(c, err.Error())
	}

	// SAML 認証成功後
	email := assertion.NameID()
	u, err := model.FindUserByEmail(email)
	if err != nil {
		return c.String(http.StatusBadRequest, "ユーザーがみつかりませんでした")
	}

	if err := session.Clear(c); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	sid, err := session.StartWithTTL(c, assertion.SessionTTL())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	session.Add(c, sid, "userId", u.ID)
	session.Add(c, sid, "success", "SAML 認証に成功しました")
	session.Activate(u.ID)

	redirect := "/"
	rs := c.FormValue("RelayState")
	if rs != "" {
		// TODO: handle RelayState
	}
	return c.Redirect(http.StatusFound, redirect)
}
