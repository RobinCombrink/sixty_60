package schema

import (
	"fmt"
	"strings"
	"time"
)

type Filter struct {
	StartDate            time.Time
	EndDate              time.Time
	SearchTerm           string
	ImportantItemFilters []ImportantItemFilter
}

type ImportantItemFilter struct {
	Name        string
	ExactSearch bool
	Blacklist   []string
}

func (f Filter) String() string {
	var importantItemFilters []string
	for _, filter := range f.ImportantItemFilters {
		importantItemFilters = append(importantItemFilters, filter.String())
	}
	return fmt.Sprintf("Filter: StartDate: %v, EndDate: %v, ImportantItemFilters: [%s]",
		f.StartDate, f.EndDate, strings.Join(importantItemFilters, ", "))
}

func (i ImportantItemFilter) String() string {
	return fmt.Sprintf("ImportantItemFilter: Name: %s, Contains: %t, Blacklist: [%s]",
		i.Name, i.ExactSearch, strings.Join(i.Blacklist, ", "))
}
