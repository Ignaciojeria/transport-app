package request

import "transport-app/app/domain"

type SearchOrdersByUniqueReferencesRequest struct {
	Commerces    []string `json:"commerces"`
	ReferenceIDs []string `json:"referenceIDs"`
}

func (r SearchOrdersByUniqueReferencesRequest) Map() domain.OrderSearchFilters {
	referenceIDs := make([]domain.ReferenceID, len(r.ReferenceIDs))
	for i, id := range r.ReferenceIDs {
		referenceIDs[i] = domain.ReferenceID(id)
	}
	return domain.OrderSearchFilters{
		ReferenceIDs: referenceIDs,
		Commerces:    r.Commerces,
	}
}
