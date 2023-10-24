package schema

import "time"

type Filter struct {
	StartDate            time.Time
	EndDate              time.Time
	ImportantItemFilters []ImportantItemFilter
}

type ImportantItemFilter struct {
	Name     string
	Contains bool
	Whitelist []string
}
