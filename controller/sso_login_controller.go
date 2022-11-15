package controller

import (
	"net/http"
	"strconv"

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
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	md := model.FindMetadtaByCompanyID(cid)
	if md == nil {
		return c.String(http.StatusNotFound, "metadata is not found")
	}

	res, err := getSAMLResponse(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	rsig := res.ResponseSignature()
	if !rsig.IsZero() {
		// TODO: validate assertion signature
	}

	asig := res.AssertionSignature()
	if !asig.IsZero() {
		// TODO: validate assertion signature
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
	// TODO: set session time limit
	session.Set(c, "userId", u.ID)
	session.Set(c, "sessionIndex", res.SessionIndex())
	session.Set(c, "success", "SAML 認証に成功しました")

	redirect := "/"
	rs := c.FormValue("RelayState")
	if rs != "" {
		redirect = rs
	}
	return c.Redirect(http.StatusFound, redirect)
}
