package schema

type DisplayTemplate interface {
	GetTemplateDisplayName() string
	GetTemplateData() interface{}
}

type DisplayInvoiceSummary struct {
	TotalSpent        string
	TotalSaved        string
	TotalItemsOrdered uint64
	TotalOrders       uint64
	TemplateName      string
	ImportantItems    map[string]DisplayImportantItem
}

type DisplayImportantItem struct {
	Names            []string
	TotalSpent       string
	MaximumUnitPrice string
	MinimumUnitPrice string
	AverageSpent     string
	AverageUnitPrice string
	TotalQuantity    uint32
	TotalSaved       string
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
	totalOrders uint64,
	importantItems map[string]DisplayImportantItem) DisplayTemplate {
	return DisplayInvoiceSummary{
		TotalSpent:        totalSpent,
		TotalSaved:        totalSaved,
		TotalItemsOrdered: totalItemsOrdered,
		TotalOrders:       totalOrders,
		ImportantItems:    importantItems,
		TemplateName:      "DisplayInvoiceSummary"}
}

type DisplayInvoiceList struct {
	Invoices     []Invoice
	TemplateName string
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
