package domain

type OrderSearchFilters struct {
	Pagination          Pagination
	Organization        Organization
	ReferenceIDs        []ReferenceID
	Packages            []Packages
	BusinessIdentifiers []BusinessIdentifiers
}
