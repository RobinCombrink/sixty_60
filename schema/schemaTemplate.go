package schema

import "strconv"

type DisplayImportantItem struct {
	Names            []string
	TotalSpent       string
	MaximumUnitPrice string
	MinimumUnitPrice string
	AverageUnitPrice string
	TotalQuantity    string
	TotalSaved       string
}

type DisplayInvoiceSummary struct {
	TotalSpent               string
	TotalSaved               string
	TotalItemsMatchingSearch string
	TotalOrders              string
}

func MakeDisplayInvoiceSummary(totalSpent string,
	totalSaved string,
	totalItemsMatchingSearch uint64,
	totalOrders uint64) DisplayInvoiceSummary {
	return DisplayInvoiceSummary{
		TotalSpent:               totalSpent,
		TotalSaved:               totalSaved,
		TotalItemsMatchingSearch: strconv.FormatUint(uint64(totalItemsMatchingSearch), 10),
		TotalOrders:              strconv.FormatUint(uint64(totalOrders), 10),
	}
}

type DisplayInvoiceList struct {
	Invoices       []Invoice
	ImportantItems map[string]DisplayImportantItem
}

func GetDisplayInvoiceList(invoices []Invoice) DisplayInvoiceList {
	return DisplayInvoiceList{
		Invoices: invoices,
		//TODO: ImportantItems
	}
}

type Page struct {
	Title string
	Url   string
}

func GetPageList() []Page {
	pages := make([]Page, 5)
	pages = append(pages, Page{
		Title: "home",
		Url:   "../home",
	})
	pages = append(pages, Page{
		Title: "invoices",
		Url:   "../invoices",
	})
	pages = append(pages, Page{
		Title: "summary",
		Url:   "../summary",
	})
	return pages
}
