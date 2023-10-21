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
	tmpl.Funcs(template.FuncMap{"divideUInt64": DivideUInt64, "divideUInt16": DivideUInt16,})
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
func DivideUInt64(param1 uint64, param2 uint64) uint64 {
	return param1 / param2
}
func DivideUInt16(param1 uint16, param2 uint16) uint16 {
	return param1 / param2
}
