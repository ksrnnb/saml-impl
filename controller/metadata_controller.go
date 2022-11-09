package controller

import (
	"net/http"

	"github.com/ksrnnb/saml/model"
	"github.com/labstack/echo/v4"
)

const defaultCompanyID = 1

// IdP のメタデータ設定ページの表示
func Metadata(c echo.Context) error {
	m := model.FindMetadtaByCompanyID(defaultCompanyID)
	return c.Render(http.StatusOK, "metadata.html", m)
}

// IdP から取得したメタデータの登録
// 既に存在する場合は上書き
func CreateMetadata(c echo.Context) error {
	m := model.NewMetadata(
		defaultCompanyID,
		c.FormValue("entityID"),
		c.FormValue("certificate"),
		c.FormValue("ssourl"),
	)
	m.Save()

	return c.Redirect(http.StatusFound, "/metadata")
}
