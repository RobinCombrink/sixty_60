package schema

import (
	"fmt"
	"strings"
	"time"
)

type DateFilter struct {
	StartDate       time.Time
	EndDate         time.Time
	// ImportantItemFilters []ImportantItemFilter
}
type ItemFilter struct {
	SearchTerm      *string
	Contains        bool
	CaseInsensitive bool
}
type ImportantItemFilter struct {
	Name        string
	ExactSearch bool
	Blacklist   []string
}

func (f DateFilter) String() string {
	var importantItemFilters []string
	// for _ , filter := range f.ImportantItemFilters {
	// importantItemFilters = append(importantItemFilters, filter.String())
	// }
	return fmt.Sprintf("Filter: StartDate: %v, EndDate: %v, ImportantItemFilters: [%s]",
		f.StartDate, f.EndDate, strings.Join(importantItemFilters, ", "))
}

func (i ItemFilter) String() string {
	return fmt.Sprintf("ItemFilter: SearchTerm: %s, Contains: %t, CaseInsensitive: %t",
		i.SearchTerm, i.Contains, i.CaseInsensitive)
}

func (i ImportantItemFilter) String() string {
	return fmt.Sprintf("ImportantItemFilter: Name: %s, Contains: %t, Blacklist: [%s]",
		i.Name, i.ExactSearch, strings.Join(i.Blacklist, ", "))
}
