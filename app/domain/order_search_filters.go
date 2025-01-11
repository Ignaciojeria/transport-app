package domain

type OrderSearchFilters struct {
	Pagination   Pagination
	Organization Organization
	ReferenceIDs []ReferenceID
	Packages     []Package
	Commerces    []string `json:"commerce"`
}
