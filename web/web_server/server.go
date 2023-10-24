package webserver

import (
	"log"
	"net/http"
	"parser60/format"
	"parser60/schema"
	"parser60/template"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const serverIp string = "127.0.0.1"
const dateLayout string = "2006-01-02"

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
	instance.GET("/home", getRoot)
	instance.GET("/invoices", getInvoices)
	instance.GET("/summary", getSummary)
	instance.POST("/summary/filter", postSummaryFilter)
}
func postSummaryFilter(c echo.Context) error {
	log.Printf("wishdbfiwjebf")
	dateStr := c.FormValue("date")
	date, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return err
	}

	displayInvoiceSummary := createInvoiceSummary(invoices, schema.Filter{Date: date})
	return c.Render(http.StatusOK, "invoicesSummary", displayInvoiceSummary)
}
func getSummary(c echo.Context) error {
	displayInvoiceSummary := createInvoiceSummary(invoices, schema.Filter{})
	return c.Render(http.StatusOK, "summary", displayInvoiceSummary)
}

func createInvoiceSummary(invoices []schema.Invoice, filter schema.Filter) schema.DisplayTemplate {
	var totalSpent uint64 = 0
	var totalSaved uint64 = 0
	var totalItemsOrdered uint64 = 0
	var totalOrders uint64 = 0

	for _, invoice := range invoices {
		if filter.Date.IsZero() || invoice.Date.After(filter.Date) {
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
	}
	// averageOrderCost =
	return schema.GetDisplayInvoiceSummary(
		format.ToRand(totalSpent),
		format.ToRand(totalSaved),
		totalItemsOrdered,
		totalOrders)
}
func getInvoices(c echo.Context) error {
	return c.Render(http.StatusOK, "invoices", schema.GetDisplayInvoiceList(invoices))
}
func getRoot(c echo.Context) error {
	return c.Render(http.StatusOK, "home", nil)
}
