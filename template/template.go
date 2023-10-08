package template

import (
	"fmt"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(writer io.Writer, name string, data interface{}, context echo.Context) error {
	fmt.Printf("data: %v\n", name)
	return t.Templates.ExecuteTemplate(writer, name, data)
}

func NewTemplateRenderer(e *echo.Echo, paths ...string) {
	tmpl := &template.Template{}
	for i := range paths {
		fmt.Printf(paths[i] + "\n")
		template.Must(tmpl.ParseGlob(paths[i]))
	}
	t := newTemplate(tmpl)
	e.Renderer = t
fmt.Printf("template: %v\n", tmpl.Lookup("index"))
}

func newTemplate(templates *template.Template) echo.Renderer {
	return &Template{
		Templates: templates,
	}
}
