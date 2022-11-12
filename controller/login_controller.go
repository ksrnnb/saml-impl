package controller

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ksrnnb/saml/model"
	"github.com/labstack/echo/v4"
)

const (
	StatusSeccess = "urn:oasis:names:tc:SAML:2.0:status:Success"
)

func ShowLogin(c echo.Context) error {
	m := model.FindMetadtaByCompanyID(defaultCompanyID)
	if m == nil {
		m = &model.Metadata{CompanyID: defaultCompanyID}
	}
	return c.Render(http.StatusOK, "login.html", m)
}

// HTTP POST Binding
func HandleSamlResponse(c echo.Context) error {
	encRes := c.FormValue("SAMLResponse")
	decRes, _ := base64.StdEncoding.DecodeString(encRes)
	res := SamlResponse{}
	if err := xml.Unmarshal(decRes, &res.Response); err != nil {
		return err
	}

	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	md := model.FindMetadtaByCompanyID(cid)
	if err := res.Validate(md); err != nil {
		return err
	}
	fmt.Printf("%+v\n", res)

	redirect := "/"
	rs := c.FormValue("RelayState")
	if rs != "" {
		redirect = rs
	}
	return c.Redirect(http.StatusFound, redirect)
}

func (r SamlResponse) Validate(md *model.Metadata) error {
	if r.Destination() != md.ACSURL() {
		return errors.New("destination is invalid")
	}

	if r.Issuer() != md.EntityID {
		return errors.New("issuer is invalid")
	}

	if r.StatusCode() != StatusSeccess {
		return errors.New("status is not success")
	}

	if r.Recipient() != md.ACSURL() {
		return errors.New("recipient is invalid")
	}

	// TODO: validate NotOnOrAfter...

	return nil
}
