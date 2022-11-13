package controller

import (
	"fmt"
	"net/http"

	"github.com/ksrnnb/saml/model"
	"github.com/ksrnnb/saml/session"
	"github.com/labstack/echo/v4"
)

const defaultCompanyID = 1

type MetadataParam struct {
	Metadata       *model.Metadata
	SuccessMessage string
}

// IdP のメタデータ設定ページの表示
func Metadata(c echo.Context) error {
	m := model.FindMetadtaByCompanyID(defaultCompanyID)
	if m == nil {
		m = &model.Metadata{CompanyID: defaultCompanyID}
	}
	sm, err := session.Get(c, "success")
	if err != nil {
		fmt.Printf("session get error: %v\n", err)
		return err
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
	m := model.NewMetadata(
		defaultCompanyID,
		c.FormValue("entityID"),
		c.FormValue("certificate"),
		c.FormValue("ssourl"),
	)
	m.Save()
	err := session.Set(c, "success", "メタデータを更新しました")
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/metadata")
}
