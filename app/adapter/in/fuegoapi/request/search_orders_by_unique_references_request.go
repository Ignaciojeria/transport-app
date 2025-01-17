package request

import "transport-app/app/domain"

type SearchOrdersByUniqueReferencesRequest struct {
	Pagination struct {
		Page int `json:"page"`
		Size int `json:"size"`
	} `json:"pagination"`
	Commerces    []string `json:"commerces"`
	ReferenceIDs []string `json:"referenceIDs"`
}

func (r SearchOrdersByUniqueReferencesRequest) Map() domain.OrderSearchFilters {
	referenceIDs := make([]domain.ReferenceID, len(r.ReferenceIDs))
	for i, id := range r.ReferenceIDs {
		referenceIDs[i] = domain.ReferenceID(id)
	}
	return domain.OrderSearchFilters{
		Pagination: domain.Pagination{
			Page: r.Pagination.Page,
			Size: r.Pagination.Size,
		},
		ReferenceIDs: referenceIDs,
		Commerces:    r.Commerces,
	}
}
