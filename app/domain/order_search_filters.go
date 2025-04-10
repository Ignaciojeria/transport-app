package domain

type OrderSearchFilters struct {
	Pagination      Pagination
	ReferenceIDs    []string
	Lpns            []string
	Commerces       []string `json:"commerce"`
	PlanReferenceID string
}
