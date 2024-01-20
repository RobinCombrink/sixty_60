package database

import (
	"parser60/emailparsing"
	"parser60/schema"
)

var invoices = make([]schema.Invoice, 50)

func addInvoice(invoice schema.Invoice) {
	invoices = append(invoices, invoice)
}

func GetInvoices() []schema.Invoice {
	return invoices
}

func AddInvoices(messageBodies []string) {
	for _, messageBody := range messageBodies {
		invoice := emailparsing.GetInvoiceFromHtml(string(messageBody))
		addInvoice(*invoice)
		invoices = append(invoices, *invoice)
	}
}
