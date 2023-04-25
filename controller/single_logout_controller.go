package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/service"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

func SingleLogout(c echo.Context) error {
	ss, err := service.NewSamlService(c.Param("company_id"))
	if err != nil {
		// TODO: エラーの場合も LogoutResponse を返す
		return err
	}

	r := c.Request()
	r.ParseForm()

	lr, err := ss.ParseLogoutRequest(r.PostForm.Get("SAMLRequest"))
	if err != nil {
		// TODO: エラーの場合も LogoutResponse を返す
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err = model.FindUserByEmail(lr.NameID())
	if err != nil {
		// TODO: エラーの場合も LogoutResponse を返す
		return c.String(http.StatusInternalServerError, "ユーザーがみつかりませんでした")
	}

	resp, err := ss.ServiceProvider.MakePostLogoutResponse(lr.ID(), "")
	if err != nil {
		// TODO: エラーの場合も LogoutResponse を返す
		return err
	}
	// ログアウトしたユーザーのセッションを無効にする
	session.Invalidate(lr.ID())
	return c.HTML(http.StatusOK, string(resp))
}
