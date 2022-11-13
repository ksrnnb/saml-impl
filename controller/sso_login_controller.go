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
	"github.com/ksrnnb/saml/session"
	"github.com/labstack/echo/v4"
)

func ShowSSOLogin(c echo.Context) error {
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	md := model.FindMetadtaByCompanyID(u.CompanyID)
	if md == nil {
		return c.Render(http.StatusOK, "ssologin.html", md)
	}
	return c.Render(http.StatusOK, "ssologin.html", md)
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

	u := model.FindUserByPersistentID(res.NameID())
	if u == nil {
		// IdP と SP の間で persistent id による紐づきがまだの場合は
		// メールアドレスで検索して紐付けを行う
		u = model.FindUserByEmail(res.Email())
		if u == nil {
			session.Set(c, "error", "IdP のメールアドレスと一致するユーザーがみつかりませんでした")
			return c.Redirect(http.StatusFound, "/login")
		}
		u.PersistentID = res.NameID()
		u.Save()
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
