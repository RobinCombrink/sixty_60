package schema

import (
	"fmt"
	"strconv"
	"time"
)

type Invoice struct {
	Items       []LineItem
	DeliveryFee uint64
	XtraSavings uint64
	Date        time.Time
}

func (invoice Invoice) CalculateItemTotals() (uint64, uint64, uint64) {
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

func (invoice Invoice) MatchesFilter(filter Filter) bool {
	return filter.StartDate.IsZero() || (invoice.Date.After(filter.StartDate) && invoice.Date.Before(filter.EndDate))
}

type LineItem struct {
	Name     string
	Quantity uint32
	Price    uint64
	Total    uint64
	Discount uint64
}

type ImportantItem struct {
	Names            map[string]Void
	MaximumUnitPrice uint64
	MinimumUnitPrice uint64
	TotalSaved       uint64
	TotalSpent       uint64
	TotalQuantity    uint32
}

func (i ImportantItem) ToDisplay() DisplayImportantItem {
	var names []string
	for name := range i.Names {
		names = append(names, name)
	}

	averageUnitPrice := "0"
	if i.TotalQuantity != 0 {
		averageUnitPrice = fmt.Sprintf("%.2f", float64(i.TotalSpent)/float64(i.TotalQuantity))
	}

	return DisplayImportantItem{
		Names:            names,
		TotalSpent:       strconv.FormatUint(i.TotalSpent, 10),
		MaximumUnitPrice: strconv.FormatUint(i.MaximumUnitPrice, 10),
		MinimumUnitPrice: strconv.FormatUint(i.MinimumUnitPrice, 10),
		AverageUnitPrice: averageUnitPrice,
		TotalQuantity:    strconv.FormatUint(uint64(i.TotalQuantity), 10),
		TotalSaved:       strconv.FormatUint(i.TotalSaved, 10),
	}
}

type Void struct{}

type Display[T any] interface {
	ToDisplay() T
}
