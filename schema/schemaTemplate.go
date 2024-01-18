package schema

type DisplayTemplate interface {
	GetTemplateDisplayName() string
	GetTemplateData() interface{}
}

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
	TemplateName      string
}

func (invoiceSummary DisplayInvoiceSummary) GetTemplateDisplayName() string {
	return invoiceSummary.TemplateName
}
func (invoiceSummary DisplayInvoiceSummary) GetTemplateData() interface{} {
	return invoiceSummary
}

func GetDisplayInvoiceSummary(totalSpent string,
	totalSaved string,
	totalItemsOrdered uint64,
	totalOrders uint64) DisplayTemplate {
	return DisplayInvoiceSummary{
		TotalSpent:        totalSpent,
		TotalSaved:        totalSaved,
		TotalItemsOrdered: totalItemsOrdered,
		TotalOrders:       totalOrders,
		TemplateName:      "DisplayInvoiceSummary"}
}

type DisplayInvoiceList struct {
	Invoices       []Invoice
	ImportantItems map[string]DisplayImportantItem
	TemplateName   string
}

func (invoiceList DisplayInvoiceList) GetTemplateDisplayName() string {
	return invoiceList.TemplateName
}

func (invoiceList DisplayInvoiceList) GetTemplateData() interface{} {
	return invoiceList
}

func GetDisplayInvoiceList(invoices []Invoice) DisplayTemplate {
	return DisplayInvoiceList{
		Invoices:     invoices,
		TemplateName: "DisplayInvoiceList"}
}

type Page struct {
	Title string
	Url   string
}

func  GetPageList() []Page{
	pages := make([]Page, 5)
	pages = append(pages, Page{
		Title: "test",
		Url:   "../home",
	})
	return pages
}
