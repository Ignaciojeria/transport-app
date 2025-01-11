package request

import "transport-app/app/domain"

type SearchOrdersRequest struct {
	Pagination struct {
		Page int `json:"page"`
		Size int `json:"size"`
	} `json:"pagination"`
	Commerces    []string `json:"commerces"`
	ReferenceIDs []string `json:"referenceIDs"`
	PackageLpns  []string `json:"packageLpns"`
}

func (r SearchOrdersRequest) Map() domain.OrderSearchFilters {
	referenceIDs := make([]domain.ReferenceID, len(r.ReferenceIDs))
	for i, id := range r.ReferenceIDs {
		referenceIDs[i] = domain.ReferenceID(id)
	}

	packages := make([]domain.Package, len(r.PackageLpns))
	for i, lpn := range r.PackageLpns {
		packages[i] = domain.Package{Lpn: lpn}
	}

	return domain.OrderSearchFilters{
		Pagination: domain.Pagination{
			Page: r.Pagination.Page,
			Size: r.Pagination.Size,
		},
		ReferenceIDs: referenceIDs,
		Packages:     packages,
		Commerces:    r.Commerces,
	}
}
