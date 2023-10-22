package webserver

import (
	"net/http"
	"parser60/format"
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
	instance.GET("summary", getSummary)
}
func getSummary(c echo.Context) error {
	displayInvoiceSummary := createInvoiceSummary(invoices)
	return c.Render(http.StatusOK, "summary", displayInvoiceSummary)
}

func createInvoiceSummary(invoices []schema.Invoice) schema.DisplayInvoiceSummay {
	var totalSpent uint64 = 0
	var totalSaved uint64 = 0
	var totalItemsOrdered uint64 = 0
	var totalOrders uint64 = 0

	for _, invoice := range invoices {
		var invoiceTotal uint64 = 0
		var invoiceSaved uint64 = 0
		var invoiceItemsOrdered uint64 = 0
		for _, lineItem := range invoice.Items {
			invoiceTotal += lineItem.Total
			invoiceSaved += lineItem.Discount
			invoiceItemsOrdered += 1
		}
		totalSpent += invoiceTotal
		totalSaved += invoiceSaved
		totalItemsOrdered += invoiceItemsOrdered
		totalOrders += 1
	}
	// averageOrderCost =
	return schema.DisplayInvoiceSummay{
		TotalSpent: format.ToRand(totalSpent),
		TotalSaved: format.ToRand(totalSaved),
		TotalItemsOrdered: totalItemsOrdered,
		TotalOrders: totalOrders}
}
func getInvoices(c echo.Context) error {
	invoiceData := make(map[string][]schema.Invoice)
	invoiceData["invoices"] = invoices
	return c.Render(http.StatusOK, "invoices", invoiceData)
}
func getRoot(c echo.Context) error {
	return c.Render(http.StatusOK, "home", nil)
}
