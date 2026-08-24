package apimodel

type ListFiltersRequest struct {
	Page           int        `form:"page" binding:"omitempty"`
	FiltersJSON    string     `form:"filtersJSON" binding:"omitempty"`
	SortOptionJSON string     `form:"sortOptionJSON" binding:"omitempty"`
	Filters        []Filter   `form:"filters" binding:"dive"`
	SortOption     SortOption `form:"sortOption"  binding:"omitempty"`
	SearchString   string     `form:"searchString"  binding:"omitempty"`
}
type Filter struct {
	Field        string   `form:"field"`
	Condition    string   `form:"condition"`
	FilterValues []string `form:"filterValues"`
}

type SortOption struct {
	SortBy   string `form:"sortBy"`
	SortType string `form:"sortType"`
}
