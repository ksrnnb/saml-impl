package controller

import (
	"net/http"

	"github.com/ksrnnb/saml/model"
	"github.com/labstack/echo/v4"
)

func ShowLogin(c echo.Context) error {
	m := model.FindMetadtaByCompanyID(defaultCompanyID)
	if m == nil {
		m = &model.Metadata{CompanyID: defaultCompanyID}
	}
	return c.Render(http.StatusOK, "login.html", m)
}

func StartSPSamlLogin(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
