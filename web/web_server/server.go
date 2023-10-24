package webserver

import (
	"fmt"
	"net/http"
	"parser60/format"
	"parser60/schema"
	"parser60/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slices"
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
	startDateStr := c.FormValue("startDate")
	endDateStr := c.FormValue("endDate")
	startDate, err := time.Parse(dateLayout, startDateStr)
	if err != nil {
		return err
	}
	endDate, err := time.Parse(dateLayout, endDateStr)
	if err != nil {
		return err
	}

	importantItemFilters := make([]schema.ImportantItemFilter, 5)
	importantItemFilters = append(importantItemFilters, schema.ImportantItemFilter{Name: "Cheese", Contains: true})
	displayInvoiceSummary := createInvoiceSummary(invoices, schema.Filter{StartDate: startDate, EndDate: endDate, ImportantItemFilters: importantItemFilters})
	return c.Render(http.StatusOK, "invoicesSummary", displayInvoiceSummary)
}
func getSummary(c echo.Context) error {
	t := time.Now()

	startDate := time.Date(t.Year(), time.January, 0, 0, 0, 0, 0, t.Location())
	endDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	importantItemFilters := make([]schema.ImportantItemFilter, 5)

	cheeseWhitelist := make([]string, 3)
	cheeseWhitelist = append(cheeseWhitelist, "Parmalat Cheddar Cheese Pack 400g")
	cheeseWhitelist = append(cheeseWhitelist, "Ladismith Cheddar Cheese Pack 800g")
	cheeseWhitelist = append(cheeseWhitelist, "Lancewood Cheddar Cheese Pack 900g")
	importantItemFilters = append(importantItemFilters, schema.ImportantItemFilter{Name: "Cheese", Contains: true, Whitelist: cheeseWhitelist})
	displayInvoiceSummary := createInvoiceSummary(invoices, schema.Filter{StartDate: startDate, EndDate: endDate, ImportantItemFilters: importantItemFilters})
	return c.Render(http.StatusOK, "summary", displayInvoiceSummary)
}

func createInvoiceSummary(invoices []schema.Invoice, filter schema.Filter) schema.DisplayTemplate {
	var totalSpent uint64 = 0
	var totalSaved uint64 = 0
	var totalItemsOrdered uint64 = 0
	var totalOrders uint64 = 0
	var importantItemsDisplay map[string]schema.DisplayImportantItem = make(map[string]schema.DisplayImportantItem, len(filter.ImportantItemFilters))
	var importantItems map[string]schema.ImportantItem = make(map[string]schema.ImportantItem, len(filter.ImportantItemFilters))
	for _, invoice := range invoices {
		if filter.StartDate.IsZero() || (invoice.Date.After(filter.StartDate) && invoice.Date.Before(filter.EndDate)) {
			var invoiceTotal uint64 = 0
			var invoiceSaved uint64 = 0
			var invoiceItemsOrdered uint64 = 0
			for _, lineItem := range invoice.Items {

				for _, importantItemFilter := range filter.ImportantItemFilters {
					if importantItemFilter.Contains {
						if strings.Contains(lineItem.Name, importantItemFilter.Name) {
							fmt.Printf("%s %v\n", lineItem.Name, lineItem.Price-lineItem.Discount)
							if slices.Contains(importantItemFilter.Whitelist, lineItem.Name) {
								importantItems[importantItemFilter.Name] = schema.ImportantItem{
									Name:             lineItem.Name,
									TotalSpent:       importantItems[importantItemFilter.Name].TotalSpent + lineItem.Total,
									TotalQuantity:    importantItems[importantItemFilter.Name].TotalQuantity + lineItem.Quantity,
									MaximumUnitPrice: getMax(importantItems[importantItemFilter.Name].MaximumUnitPrice, lineItem.Price-lineItem.Discount),
									MinimumUnitPrice: getMin(importantItems[importantItemFilter.Name].MinimumUnitPrice, lineItem.Price-lineItem.Discount),
									TotalSaved:       importantItems[importantItemFilter.Name].TotalSaved + (lineItem.Discount * uint64(lineItem.Quantity)),
								}
							}
						}
					}
					// else {
					// 	if lineItem.Name == importantItemFilter.Name {

					// 	}
					// }
				}

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
	for key, importantItem := range importantItems {
		importantItemsDisplay[key] = schema.DisplayImportantItem{
			Name:             importantItem.Name,
			TotalSpent:       format.ToRand(importantItem.TotalSpent),
			TotalQuantity:    importantItem.TotalQuantity,
			TotalSaved:       format.ToRand(importantItem.TotalSaved),
			MaximumUnitPrice: format.ToRand(importantItem.MaximumUnitPrice),
			MinimumUnitPrice: format.ToRand(importantItem.MinimumUnitPrice),
			AverageSpent:     format.ToRand(importantItem.TotalSpent / uint64(importantItem.TotalQuantity)),
			// AverageUnitPrice: format.ToRand(importantItem.),
		}
	}

	return schema.GetDisplayInvoiceSummary(
		format.ToRand(totalSpent),
		format.ToRand(totalSaved),
		totalItemsOrdered,
		totalOrders,
		importantItemsDisplay)
}

func getMax(num1 uint64, num2 uint64) uint64 {
	if num1 > num2 {
		return num1
	} else {
		return num2
	}
}
func getMin(num1 uint64, num2 uint64) uint64 {
	if num1 == 0 {
		return num2
	} else if num2 == 0 {
		return num1
	}
	if num1 < num2 {
		return num1
	} else {
		return num2
	}
}

func getInvoices(c echo.Context) error {
	return c.Render(http.StatusOK, "invoices", schema.GetDisplayInvoiceList(invoices))
}
func getRoot(c echo.Context) error {
	return c.Render(http.StatusOK, "home", nil)
}
