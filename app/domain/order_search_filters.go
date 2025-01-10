package domain

type OrderSearchFilters struct {
	Pagination   Pagination
	Organization Organization
	ReferenceIDs []ReferenceID
	Packages     []Packages
	Commerces    []string `json:"commerce"`
}
