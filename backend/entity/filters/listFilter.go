package filters

type ListFilter struct {
	Page         int
	Filters      []Filter
	SortOption   SortOption
	SearchString string
}
type Filter struct {
	Field        string
	Condition    string
	FilterValues []string
}

type SortOption struct {
	SortBy   string
	SortType string
}
