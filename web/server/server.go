package server

import (
	"context"
	"fmt"
	"net/http"
	"parser60/database"
	"parser60/format"
	"parser60/schema"
	"parser60/template"
	genericTemplates "parser60/web/templates"
	authTemplates "parser60/web/templates/auth"
	invoiceTemplates "parser60/web/templates/invoices"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const serverIp string = "127.0.0.1"
const dateLayout string = "2006-01-02"

var serverContext echo.Context
var blacklist []string

func SetupOAuthServer(result chan<- string) {
	println("Setting up OAuth Server")
	echoInstance := echo.New()

	echoInstance.Pre(middleware.RemoveTrailingSlash())

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())

	setupAuthRoutes(echoInstance, result)
	echoInstance.Logger.Fatal(echoInstance.Start(serverIp + ":42069"))
}

func setupAuthRoutes(authInstance *echo.Echo, result chan<- string) {
	authInstance.GET("/", func(c echo.Context) error {
		return getAuthRoot(c, result)
	})
}

func getAuthRoot(c echo.Context, result chan<- string) error {
	code := c.QueryParam("code")
	go func() {
		result <- code
	}()

	return render(c, http.StatusOK, authTemplates.CreateAuthPage())
}

/*
Starts echo HTTP server and blocks the current thread
*/
func SetupHttpServer() {
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
	// instance.GET("/summary/important", getImportantItems)
	instance.POST("/important/filter", postImportantItemsFilter)
	instance.DELETE("/summary/important/delete/{importantItemName}", deleteImportantItem)
}

func getRoot(c echo.Context) error {
	return render(c, http.StatusOK, genericTemplates.CreateHomePage(schema.GetPageList()))
}
func getInvoices(c echo.Context) error {
	// serverContext.Logger().Printf("Important Items Len: %v", len(importantItems))

	// var importantItemsDisplay map[string]schema.DisplayImportantItem = createDisplayImportantItems(importantItems)
	return render(c, http.StatusOK, invoiceTemplates.CreateInvoicesListPage(schema.GetDisplayInvoiceList(database.GetInvoices()).Invoices))
}
func getSummary(c echo.Context) error {
	serverContext = c
	t := time.Now()

	//TODO: Timezones :)
	//The hour is "14" because if it's 0 then it will drop back a day because of timezones
	startDate := time.Date(t.Year()-1, time.January, 14, 0, 0, 0, 0, t.Location())
	endDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	importantItemFilters := make([]schema.ImportantItemFilter, 5)

	// cheeseWhitelist := make([]string, 3)
	// cheeseWhitelist = append(cheeseWhitelist, "Parmalat Cheddar Cheese Pack 400g")
	// cheeseWhitelist = append(cheeseWhitelist, "Ladismith Cheddar Cheese Pack 800g")
	// cheeseWhitelist = append(cheeseWhitelist, "Lancewood Cheddar Cheese Pack 900g")
	// importantItemFilters = append(importantItemFilters, schema.ImportantItemFilter{Name: "Cheese", Contains: true})

	dateFilter := schema.Filter{StartDate: startDate, EndDate: endDate, ImportantItemFilters: importantItemFilters}

	fmt.Println("\n", dateFilter, "\n")

	displayInvoiceSummary, displayImportantItems := createInvoiceSummary(database.GetInvoices(), dateFilter)
	return render(c, http.StatusOK, invoiceTemplates.CreateInvoiceSummaryPage(startDate.Format(dateLayout), endDate.Format(dateLayout), displayInvoiceSummary, displayImportantItems))
}

// func getImportantItems(c echo.Context) error {

//		importantItems:= make([] schema.ImportantItem, 1)
//		importantItems = append(importantItems, schema.ImportantItem{})
//		displayImportantItems := createDisplayImportantItems();
//		return render(c, http.StatusOK, createImportantItemsPage())
//	}
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

	fmt.Println("ContainsStr:", containsStr)

	var contains bool
	if containsStr == "on" {
		contains = true
	} else {
		contains = false
	}
	importantItemFilters := make([]schema.ImportantItemFilter, 0)
	importantItemFilters = append(importantItemFilters, schema.ImportantItemFilter{Name: searchText, ExactSearch: contains})

	dateFilter := schema.Filter{
		StartDate:            startDate,
		EndDate:              endDate,
		ImportantItemFilters: importantItemFilters,
	}

	fmt.Println("\n", dateFilter, "\n")

	displayInvoiceSummary, _ := createInvoiceSummary(database.GetInvoices(), dateFilter)
	return render(c, http.StatusOK, invoiceTemplates.CreateInvoicesSummarySection(displayInvoiceSummary))
}

func postImportantItemsFilter(c echo.Context) error {
	return c.Render(http.StatusOK, "Not Implemented yet", nil)
}

func deleteImportantItem(c echo.Context) error {
	return c.Render(http.StatusOK, "empty", nil)
}

func updateImportantItems(importantItems map[string]schema.ImportantItem, lineItem schema.LineItem, importantItemFilter schema.ImportantItemFilter) map[string]schema.ImportantItem {
	if importantItemFilter.ExactSearch {
		if lineItem.Name == importantItemFilter.Name {
			checkBlacklistAndFilter(importantItemFilter, lineItem, importantItems)
		}
	} else {
		if strings.Contains(lineItem.Name, importantItemFilter.Name) {
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
	// serverContext.Logger().Printf("%s %v\n", lineItem.Name, lineItem.Price-lineItem.Discount)
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
			TotalQuantity:    strconv.FormatUint(uint64(importantItem.TotalQuantity), 10),
			TotalSaved:       format.ToRand(importantItem.TotalSaved),
			MaximumUnitPrice: format.ToRand(importantItem.MaximumUnitPrice),
			MinimumUnitPrice: format.ToRand(importantItem.MinimumUnitPrice),
			AverageUnitPrice: format.ToRand(importantItem.TotalSpent / uint64(importantItem.TotalQuantity)),
		}
	}
	return importantItemsDisplay
}

func createInvoiceSummary(invoices []schema.Invoice, filter schema.Filter) (schema.DisplayInvoiceSummary, []schema.DisplayImportantItem) {
	var totalSpent uint64 = 0
	var totalSaved uint64 = 0
	var totalItemsOrdered uint64 = 0
	var totalOrders uint64 = 0

	var importantItems map[string]schema.ImportantItem = make(map[string]schema.ImportantItem, len(filter.ImportantItemFilters))

	for _, invoice := range invoices {
		if invoice.MatchesFilter(filter) {

			invoiceTotal, invoiceSaved, invoiceItemsOrdered := invoice.CalculateItemTotals()

			for _, lineItem := range invoice.Items {
				for _, importantItemFilter := range filter.ImportantItemFilters {
					importantItems = updateImportantItems(importantItems, lineItem, importantItemFilter)
				}
			}

			totalSpent += invoiceTotal
			totalSaved += invoiceSaved
			totalItemsOrdered += invoiceItemsOrdered
			totalOrders += 1
		}
	}

	var importantItemsDisplays []schema.DisplayImportantItem = make([]schema.DisplayImportantItem, len(filter.ImportantItemFilters))

	for _, importantItem := range importantItems {
		importantItemsDisplays = append(importantItemsDisplays, importantItem.ToDisplay())
	}

	displayInvoiceSummary := schema.DisplayInvoiceSummary{
		TotalSpent:        format.ToRand(totalSpent),
		TotalSaved:        format.ToRand(totalSaved),
		TotalItemsOrdered: totalItemsOrdered,
		TotalOrders:       totalOrders}
	return displayInvoiceSummary, importantItemsDisplays
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

func render(c echo.Context, status int, t templ.Component) error {
	c.Response().Writer.WriteHeader(status)

	err := t.Render(context.Background(), c.Response().Writer)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to render response template")
	}

	return nil
}
