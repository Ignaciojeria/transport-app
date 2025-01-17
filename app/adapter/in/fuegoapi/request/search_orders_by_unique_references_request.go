package request

import "transport-app/app/domain"

type SearchOrdersByUniqueReferencesRequest struct {
	Commerces    []string `json:"commerces"`
	ReferenceIDs []string `json:"referenceIDs"`
}

func (r SearchOrdersByUniqueReferencesRequest) Map() domain.OrderSearchFilters {
	return domain.OrderSearchFilters{
		ReferenceIDs: r.ReferenceIDs,
		Commerces:    r.Commerces,
	}
}
