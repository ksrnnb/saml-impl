package main

import (
	"io"
	"text/template"

	"github.com/ksrnnb/saml-impl/route"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
	e := echo.New()
	e.Renderer = t

	route.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":3000"))
}
