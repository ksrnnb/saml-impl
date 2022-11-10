package controller

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
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

func HandleSamlResponse(c echo.Context) error {
	encRes := c.FormValue("SAMLResponse")
	decRes, _ := base64.StdEncoding.DecodeString(encRes)
	res := SamlResponse{}

	xml.Unmarshal(decRes, &res.Response)

	fmt.Printf("%+v\n", res)
	return c.Redirect(http.StatusFound, "/")
}

func (r SamlResponse) Validate() error {
	return nil
}
