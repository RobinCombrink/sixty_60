package webserver

import (
	"net/http"
	"parser60/format"
	"parser60/schema"
	"parser60/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const serverIp string = "127.0.0.1"
const dateLayout string = "2006-01-02"

var invoices []schema.Invoice
var serverContext echo.Context
var blacklist []string

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
	instance.DELETE("/summary/important_items/delete/{importantItemName}", deleteImportantItem)
}
func deleteImportantItem(c echo.Context) error {
	return c.Render(http.StatusOK, "empty", nil)
}
func postSummaryFilter(c echo.Context) error {
	serverContext = c
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
	searchText := c.FormValue("searchText")
	containsStr := c.FormValue("searchContains")
	var contains bool
	if containsStr == "true" {
		contains = true
	} else {
		contains = false
	}
	importantItemFilters := make([]schema.ImportantItemFilter, 1)
	importantItemFilters = append(importantItemFilters, schema.ImportantItemFilter{Name: searchText, Contains: contains})
	displayInvoiceSummary := createInvoiceSummary(invoices, schema.Filter{StartDate: startDate, EndDate: endDate, ImportantItemFilters: importantItemFilters})
	return c.Render(http.StatusOK, "invoicesSummary", displayInvoiceSummary)
}
func getSummary(c echo.Context) error {
	serverContext = c
	t := time.Now()

	startDate := time.Date(t.Year(), time.January, 0, 0, 0, 0, 0, t.Location())
	endDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	importantItemFilters := make([]schema.ImportantItemFilter, 5)

	// cheeseWhitelist := make([]string, 3)
	// cheeseWhitelist = append(cheeseWhitelist, "Parmalat Cheddar Cheese Pack 400g")
	// cheeseWhitelist = append(cheeseWhitelist, "Ladismith Cheddar Cheese Pack 800g")
	// cheeseWhitelist = append(cheeseWhitelist, "Lancewood Cheddar Cheese Pack 900g")
	// importantItemFilters = append(importantItemFilters, schema.ImportantItemFilter{Name: "Cheese", Contains: true})
	displayInvoiceSummary := createInvoiceSummary(invoices, schema.Filter{StartDate: startDate, EndDate: endDate, ImportantItemFilters: importantItemFilters})
	return c.Render(http.StatusOK, "summary", displayInvoiceSummary)
}

func calculateInvoice(invoice schema.Invoice, filter schema.Filter) (uint64, uint64, uint64) {
	var invoiceTotal uint64 = 0
	var invoiceSaved uint64 = 0
	var invoiceItemsOrdered uint64 = 0
	for _, lineItem := range invoice.Items {
		invoiceTotal += lineItem.Total
		invoiceSaved += lineItem.Discount
		invoiceItemsOrdered += 1
	}
	return invoiceTotal, invoiceSaved, invoiceItemsOrdered
}

func updateImportantItems(importantItems map[string]schema.ImportantItem, lineItem schema.LineItem, importantItemFilter schema.ImportantItemFilter) map[string]schema.ImportantItem {
	if importantItemFilter.Contains {
		if strings.Contains(lineItem.Name, importantItemFilter.Name) {
			checkBlacklistAndFilter(importantItemFilter, lineItem, importantItems)
		}
	} else {
		if lineItem.Name == importantItemFilter.Name {
			checkBlacklistAndFilter(importantItemFilter, lineItem, importantItems)
		}
	}
	return importantItems
}

func checkBlacklistAndFilter(importantItemFilter schema.ImportantItemFilter, lineItem schema.LineItem, importantItems map[string]schema.ImportantItem) {
	if importantItemFilter.Blacklist != nil {
		if !slices.Contains(importantItemFilter.Blacklist, lineItem.Name) {
			filter(lineItem, importantItems, importantItemFilter)
		}
	} else {
		filter(lineItem, importantItems, importantItemFilter)
	}
}

func filter(lineItem schema.LineItem, importantItems map[string]schema.ImportantItem, importantItemFilter schema.ImportantItemFilter) {
	serverContext.Logger().Printf("%s %v\n", lineItem.Name, lineItem.Price-lineItem.Discount)
	names := importantItems[importantItemFilter.Name].Names
	if names == nil {
		names = make(map[string]schema.Void)
	}
	names[lineItem.Name] = schema.Void{}
	importantItems[importantItemFilter.Name] = schema.ImportantItem{
		Names:            names,
		TotalSaved:       importantItems[importantItemFilter.Name].TotalSaved + (lineItem.Discount * uint64(lineItem.Quantity)),
		TotalSpent:       importantItems[importantItemFilter.Name].TotalSpent + lineItem.Total,
		TotalQuantity:    importantItems[importantItemFilter.Name].TotalQuantity + lineItem.Quantity,
		MaximumUnitPrice: getMax(importantItems[importantItemFilter.Name].MaximumUnitPrice, lineItem.Price-lineItem.Discount),
		MinimumUnitPrice: getMin(importantItems[importantItemFilter.Name].MinimumUnitPrice, lineItem.Price-lineItem.Discount),
	}
}

func createDisplayImportantItems(importantItems map[string]schema.ImportantItem) map[string]schema.DisplayImportantItem {
	var importantItemsDisplay map[string]schema.DisplayImportantItem = make(map[string]schema.DisplayImportantItem)
	for key, importantItem := range importantItems {
		importantItemsDisplay[key] = schema.DisplayImportantItem{
			Names:            maps.Keys(importantItem.Names),
			TotalSpent:       format.ToRand(importantItem.TotalSpent),
			TotalQuantity:    importantItem.TotalQuantity,
			TotalSaved:       format.ToRand(importantItem.TotalSaved),
			MaximumUnitPrice: format.ToRand(importantItem.MaximumUnitPrice),
			MinimumUnitPrice: format.ToRand(importantItem.MinimumUnitPrice),
			AverageUnitPrice: format.ToRand(importantItem.TotalSpent / uint64(importantItem.TotalQuantity)),
		}
	}
	return importantItemsDisplay
}

func createInvoiceSummary(invoices []schema.Invoice, filter schema.Filter) schema.DisplayTemplate {
	var totalSpent uint64 = 0
	var totalSaved uint64 = 0
	var totalItemsOrdered uint64 = 0
	var totalOrders uint64 = 0

	var importantItems map[string]schema.ImportantItem = make(map[string]schema.ImportantItem, len(filter.ImportantItemFilters))

	for _, invoice := range invoices {
		if filter.StartDate.IsZero() || (invoice.Date.After(filter.StartDate) && invoice.Date.Before(filter.EndDate)) {

			invoiceTotal, invoiceSaved, invoiceItemsOrdered := calculateInvoice(invoice, filter)

			for _, lineItem := range invoice.Items {
				for _, importantItemFilter := range filter.ImportantItemFilters {
					impo rtantItems = updateImportantItems(importantItems, lineItem, importantItemFilter)
				}
			}

			totalSpent += invoiceTotal
			totalSaved += invoiceSaved
			totalItemsOrdered += invoiceItemsOrdered
			totalOrders += 1
		}
	}

	return schema.GetDisplayInvoiceSummary(
		format.ToRand(totalSpent),
		format.ToRand(totalSaved),
		totalItemsOrdered,
		totalOrders,)
}

func getMax(num1, num2 uint64) uint64 {
	if num1 > num2 {
		return num1
	}
	return num2
}

func getMin(num1, num2 uint64) uint64 {
	if num1 == 0 || num1 > num2 {
		return num2
	}
	return num1
}

func getInvoices(c echo.Context) error {
	serverContext.Logger().Printf("Important Items Len: %v", len(importantItems))

	var importantItemsDisplay map[string]schema.DisplayImportantItem = createDisplayImportantItems(importantItems)
	return c.Render(http.StatusOK, "invoices", schema.GetDisplayInvoiceList(invoices))
}
func getRoot(c echo.Context) error {
	return c.Render(http.StatusOK, "home", nil)
}
