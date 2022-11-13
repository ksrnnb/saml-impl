package controller

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		return c.String(http.StatusBadRequest, err.Error())
	}

	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	md := model.FindMetadtaByCompanyID(cid)
	if md == nil {
		return c.String(http.StatusNotFound, "metadata is not found")
	}

	if err := res.Validate(md); err != nil {

		return c.String(http.StatusBadRequest, err.Error())
	}

	// TODO: issue session
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

	cnb, err := r.ConditionNotBefore()
	if err != nil {
		return fmt.Errorf("condition NotBefore: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	cnooa, err := r.ConditionNotOnOrAfter()
	if err != nil {
		return fmt.Errorf("parse error: condition NotOnOrAfter: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	snooa, err := r.SubjectNotOnOrAfter()
	if err != nil {
		return fmt.Errorf("parse error: subject NotOnOrAfter: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	_, err = r.SessionNotOnOrAfter()
	if err != nil {
		return fmt.Errorf("parse error: session NotOnOrAfter: %v, now: %v", err, time.Now().Format(time.RFC3339))
	}
	now := time.Now()

	if now.Before(cnb) {
		return fmt.Errorf("condition NotBefore: %v, now: %v", cnb, time.Now().Format(time.RFC3339))
	}
	if !now.Before(cnooa) {
		return fmt.Errorf("condition NotOnOrAfter: %v, now: %v", cnb, time.Now().Format(time.RFC3339))
	}
	if !now.Before(snooa) {
		return fmt.Errorf("subject NotOnOrAfter: %v, now: %v", cnb, time.Now().Format(time.RFC3339))
	}
	return nil
}
