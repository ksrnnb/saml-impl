package controller

import (
	"io/ioutil"
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
	idpMD, err := model.FindMetadtaByCompanyID(u.CompanyID)
	if idpMD == nil {
		idpMD = &model.IdPMetadata{CompanyID: u.CompanyID}
	}
	s := samlSPService(u.CompanyID)
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

// IdP からアップロードしたメタデータを Parse する
// JSON を返す
func ParseMetadata(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "cannot read request body"})
	}
	s := samlIdPService()
	m, err := s.Parse(body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, m)
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
		u.CompanyID,
		c.FormValue("entityID"),
		c.FormValue("certificate"),
		c.FormValue("ssourl"),
		c.FormValue("slourl"),
	)
	m.Save()
	err = session.Set(c, "success", "メタデータを更新しました")
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/metadata")
}

func samlSPService(cid string) service.SamlSPService {
	return service.NewSamlSPService(cid)
}

func samlIdPService() service.SamlIdPService {
	return service.NewSamlIdPService()
}
