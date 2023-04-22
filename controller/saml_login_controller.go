package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/service"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

const NotFoundCompany = "会社がみつかりませんでした。会社 ID を確認してください。"

type ShowSAMLLoginArg struct {
	Message   string
	Companies []*model.Company
}

func ShowSAMLLogin(c echo.Context) error {
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}
	errMsg, err := session.Get(c, "error")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	companies, err := model.ListAllCompanies()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	arg := ShowSAMLLoginArg{
		Message:   errMsg,
		Companies: companies,
	}
	return c.Render(http.StatusOK, "saml_login.html", arg)
}

// SAMLLogin starts to SP-initiated authentication.
func SAMLLogin(c echo.Context) error {
	uid, err := notAuthenticate(c)
	if err != nil || uid != "" {
		return err
	}

	company, err := model.FindCompany(c.FormValue("company_id"))
	if err != nil {
		return errorRedirectToSAMLLogin(c, err.Error())
	}
	if company == nil {
		return errorRedirectToSAMLLogin(c, NotFoundCompany)
	}

	ss, err := service.NewSamlService(company.ID)
	if err != nil {
		return err
	}
	u, err := ss.MakeAuthnRequestURL("")
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, u.String())
}
