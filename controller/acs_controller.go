package controller

import (
	"net/http"

	"github.com/crewjam/saml/samlsp"
	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

// HTTP POST Binding
func ConsumeAssertion(c echo.Context) error {
	md, err := model.FindMetadtaByCompanyID(c.Param("company_id"))
	if err != nil {
		return err
	}
	if md == nil {
		return c.String(http.StatusNotFound, "metadata is not found")
	}

	ss := samlSPService(md.CompanyID)
	is := samlIdPService()
	ied, err := is.BuildIdPEntityDescriptor(md)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	samlsp, _ := samlsp.New(samlsp.Options{
		EntityID:          ss.SPEntityID().String(),
		AllowIDPInitiated: true,
		IDPMetadata:       ied,
	})
	samlsp.ServiceProvider.AcsURL = *ss.ACSURL()
	samlsp.ServiceProvider.SloURL = *ss.SLOURL()

	r := c.Request()
	r.ParseForm()

	// TODO: handle SP-initiated
	possibleRequestIDs := []string{}
	if samlsp.ServiceProvider.AllowIDPInitiated {
		possibleRequestIDs = append(possibleRequestIDs, "")
	}

	assertion, err := samlsp.ServiceProvider.ParseResponse(r, possibleRequestIDs)
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

	redirect := "/"
	rs := c.FormValue("RelayState")
	if rs != "" {
		// TODO: handle RelayState
	}
	return c.Redirect(http.StatusFound, redirect)
}
