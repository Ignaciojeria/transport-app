package domain

type OrderSearchFilters struct {
	Pagination   Pagination
	Organization Organization
	ReferenceIDs []string
	Lpns         []string
	Commerces    []string `json:"commerce"`
}
