package webserver

import (
	"net/http"
	"parser60/schema"
	"parser60/template"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const serverIp = "127.0.0.1"

var invoices []schema.Invoice

/*
Starts echo HTTP server and blocks the current thread
*/
func SetupHttpServer(incomingInvoices []schema.Invoice) {
	invoices = incomingInvoices
	instance := echo.New()
	instance.Pre(middleware.RemoveTrailingSlash())

	instance.Use(middleware.Logger())
	instance.Use(middleware.Recover())

	template.NewTemplateRenderer(instance, filepath.Join("web", "templates", "*.html"))

	setupRoutes(instance)
	instance.Logger.Fatal(instance.Start(serverIp + ":42069"))
}

func setupRoutes(instance *echo.Echo) {
	instance.GET("/", getRoot)
	instance.GET("home", getRoot)
	instance.GET("invoices", getInvoices)
}
func getInvoices(c echo.Context) error {
	invoiceData := make(map[string][]schema.Invoice)
	invoiceData["invoices"] = invoices
	return c.Render(http.StatusOK, "invoices", invoiceData)
}
func getRoot(c echo.Context) error {
	return c.Render(http.StatusOK, "home", nil)
}