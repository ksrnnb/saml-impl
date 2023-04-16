package controller

import (
	"net/http"

	"github.com/ksrnnb/saml-impl/model"
	"github.com/ksrnnb/saml-impl/service"
	"github.com/ksrnnb/saml-impl/session"
	"github.com/labstack/echo/v4"
)

type MetadataParam struct {
	IdPMetadata    *model.IdPMetadata
	SPMetadata     *model.SPMetadata
	SuccessMessage string
}

// IdP のメタデータ設定ページの表示
func Metadata(c echo.Context) error {
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	idpMD := model.FindMetadtaByCompanyID(u.Company.ID)
	if idpMD == nil {
		idpMD = &model.IdPMetadata{CompanyID: u.Company.ID}
	}
	s := samlService(u.Company.ID)
	spMD := s.SPMetadata()

	sm, err := session.Get(c, "success")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(
		http.StatusOK,
		"metadata.html",
		MetadataParam{
			IdPMetadata:    idpMD,
			SPMetadata:     spMD,
			SuccessMessage: sm,
		})
}

func samlService(cid string) service.SamlService {
	return service.NewSamlService(cid)
}

// IdP から取得したメタデータの登録
// 既に存在する場合は上書き
func CreateMetadata(c echo.Context) error {
	// NOTE: need to protect from csrf
	u, err := authenticate(c)
	if err != nil || u == nil {
		return err
	}
	m := model.NewIdPMetadata(
		u.Company.ID,
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
