package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// IdP のメタデータ設定ページの表示
func Metadata(c echo.Context) error {
	return c.Render(http.StatusOK, "metadata.html", nil)
}

// IdP から取得したメタデータの登録
func CreateMetadata(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
