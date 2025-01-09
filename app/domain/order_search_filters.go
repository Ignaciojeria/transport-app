package domain

type OrderSearchFilters struct {
	Pagination struct {
		Page int
		Size int
	}
	Organization        Organization
	ReferenceIDs        []ReferenceID
	Packages            []Packages
	BusinessIdentifiers []BusinessIdentifiers
}
