package schema

type DisplayImportantItem struct {
	Names            []string
	TotalSpent       string
	MaximumUnitPrice string
	MinimumUnitPrice string
	AverageUnitPrice string
	TotalQuantity    uint32
	TotalSaved       string
}

type DisplayInvoiceSummary struct {
	TotalSpent        string
	TotalSaved        string
	TotalItemsOrdered uint64
	TotalOrders       uint64
}

func GetDisplayInvoiceSummary(totalSpent string,
	totalSaved string,
	totalItemsOrdered uint64,
	totalOrders uint64) DisplayInvoiceSummary {
	return DisplayInvoiceSummary{
		TotalSpent:        totalSpent,
		TotalSaved:        totalSaved,
		TotalItemsOrdered: totalItemsOrdered,
		TotalOrders:       totalOrders,
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
	return pages
}
