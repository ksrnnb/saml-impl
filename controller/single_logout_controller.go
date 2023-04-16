package controller

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

func SingleLogout(c echo.Context) error {
	md, err := model.FindMetadtaByCompanyID(c.Param("company_id"))
	if err != nil {
		// TODO: LogoutResponse を返す
		// 他のエラーハンドリングも同じく
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

	rawRequestBuf, err := base64.StdEncoding.DecodeString(r.PostForm.Get("SAMLRequest"))
	if err != nil {

	}
	var logoutRequest saml.LogoutRequest
	if err := xml.Unmarshal(rawRequestBuf, &logoutRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	email := logoutRequest.NameID.Value
	u, err := model.FindUserByEmail(email)
	if err != nil {
		return c.String(http.StatusInternalServerError, "ユーザーがみつかりませんでした")
	}

	fmt.Println(u)
	resp, err := samlsp.ServiceProvider.MakePostLogoutResponse(samlsp.ServiceProvider.GetSLOBindingLocation(saml.HTTPPostBinding), logoutRequest.ID)
	if err != nil {
		return err
	}
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	// TODO: 特定のユーザーのセッションを削除する
	session.Clear(c)
	return c.HTML(http.StatusOK, string(resp))
}
