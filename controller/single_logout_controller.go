package controller

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/crewjam/saml"
	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/service"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

func SingleLogout(c echo.Context) error {
	ss, err := service.NewSamlService(c.Param("company_id"))
	if err != nil {
		// TODO: LogoutResponse を返す
		// 他のエラーハンドリングも同じく
		return err
	}

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

	resp, err := ss.ServiceProvider.MakePostLogoutResponse(logoutRequest.ID, "")
	if err != nil {
		return err
	}
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	// ログアウトしたユーザーのセッションを無効にする
	session.Invalidate(u.ID)
	return c.HTML(http.StatusOK, string(resp))
}
