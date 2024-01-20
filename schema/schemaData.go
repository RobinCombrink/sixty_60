package schema

import (
	"parser60/format"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

type Invoice struct {
	Items       []LineItem
	DeliveryFee uint64
	XtraSavings uint64
	Date        time.Time
}

func (invoice Invoice) CalculateItemTotals(itemFilter ItemFilter) (invoiceTotal uint64, invoiceSaved uint64, invoiceItemsOrdered uint64) {
	for _, lineItem := range invoice.Items {
		if lineItem.MatchesFilter(itemFilter) {
			invoiceTotal += lineItem.Total
			invoiceSaved += lineItem.Discount
			invoiceItemsOrdered += 1
		}
	}
	return invoiceTotal, invoiceSaved, invoiceItemsOrdered
}

func (invoice Invoice) MatchesFilter(filter DateFilter) bool {
	return filter.StartDate.IsZero() || (invoice.Date.After(filter.StartDate) && invoice.Date.Before(filter.EndDate))
}

type LineItem struct {
	Name     string
	Quantity uint32
	Price    uint64
	Total    uint64
	Discount uint64
}

func (lineItem LineItem) MatchesFilter(filter ItemFilter) bool {

	if filter.SearchTerm == nil {
		return true
	}

	name := lineItem.Name
	searchTerm := *filter.SearchTerm

	if filter.CaseInsensitive {
		name = strings.ToLower(name)
		searchTerm = strings.ToLower(searchTerm)
	}

	if filter.Contains {
		return strings.Contains(name, searchTerm)
	}

	return name == searchTerm
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
	//  	averageUnitPrice := "0"
	// if i.TotalQuantity != 0 {
	// averageUnitPrice = fmt.Sprintf("%.2f", float64(i.TotalSpent)/float64(i.TotalQuantity))
	//  }

	return DisplayImportantItem{
		Names:            maps.Keys(i.Names),
		TotalSpent:       format.ToRand(i.TotalSpent),
		TotalQuantity:    strconv.FormatUint(uint64(i.TotalQuantity), 10),
		TotalSaved:       format.ToRand(i.TotalSaved),
		MaximumUnitPrice: format.ToRand(i.MaximumUnitPrice),
		MinimumUnitPrice: format.ToRand(i.MinimumUnitPrice),
		AverageUnitPrice: format.ToRand(i.TotalSpent / uint64(i.TotalQuantity)),
	}
}

type Void struct{}

type Display[T any] interface {
	ToDisplay() T
}

type Filterable interface {
	MatchesFilter() bool
}
