package controller

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleLogoutRequest(c echo.Context) error {
	encRes := c.FormValue("SAMLRequest")
	decRes, _ := base64.StdEncoding.DecodeString(encRes)

	req := LogoutRequest{}
	if err := xml.Unmarshal(decRes, &req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("%+v\n", req)
	return c.String(http.StatusOK, "SLO OK")
}
