package route

import (
	"github.com/ksrnnb/saml-impl/controller"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", controller.Home)

	e.GET("/metadata", controller.Metadata)
	e.POST("/metadata", controller.CreateMetadata)
	e.POST("/metadata/parse", controller.ParseMetadata)

	e.GET("/login", controller.ShowLogin)
	e.GET("/login/saml", controller.ShowSAMLLogin)
	e.POST("/login", controller.Login)
	e.POST("/login/saml", controller.SAMLLogin)

	e.POST("/logout", controller.Logout)

	e.POST("/acs/:company_id", controller.ConsumeAssertion)
	e.POST("/slo/:company_id", controller.SingleLogout)
}
