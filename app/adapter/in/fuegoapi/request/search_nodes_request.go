package request

import "transport-app/app/domain"

type SearchNodesRequest struct {
	Pagination struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	} `json:"pagination"`
}

func (s SearchNodesRequest) Map() domain.NodeSearchFilters {
	return domain.NodeSearchFilters{
		Pagination: domain.Pagination{
			Page: s.Pagination.Page,
			Size: s.Pagination.PageSize,
		},
	}
}
