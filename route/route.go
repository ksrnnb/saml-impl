package route

import (
	"github.com/ksrnnb/saml/controller"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", controller.Home)

	e.GET("/metadata", controller.Metadata)
	e.POST("/metadata", controller.CreateMetadata)

	e.GET("/login", controller.ShowLogin)
	// e.GET("/login/saml", controller.StartSPSamlLogin)
	e.POST("/login/saml/companies/:id", controller.HandleSamlResponse)
}
