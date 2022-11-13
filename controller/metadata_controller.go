package controller

import (
	"net/http"

	"github.com/ksrnnb/saml/model"
	"github.com/ksrnnb/saml/session"
	"github.com/labstack/echo/v4"
)

type MetadataParam struct {
	Metadata       *model.Metadata
	SuccessMessage string
}

// IdP のメタデータ設定ページの表示
func Metadata(c echo.Context) error {
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	m := model.FindMetadtaByCompanyID(u.CompanyID)
	if m == nil {
		m = &model.Metadata{CompanyID: u.CompanyID}
	}
	sm, err := session.Get(c, "success")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(
		http.StatusOK,
		"metadata.html",
		MetadataParam{
			Metadata:       m,
			SuccessMessage: sm,
		})
}

// IdP から取得したメタデータの登録
// 既に存在する場合は上書き
func CreateMetadata(c echo.Context) error {
	// NOTE: need to protect from csrf
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	m := model.NewMetadata(
		u.CompanyID,
		c.FormValue("entityID"),
		c.FormValue("certificate"),
		c.FormValue("ssourl"),
	)
	m.Save()
	err = session.Set(c, "success", "メタデータを更新しました")
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/metadata")
}
