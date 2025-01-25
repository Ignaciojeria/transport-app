package request

import "transport-app/app/domain"

type SearchCarriersRequest struct {
	Pagination struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	} `json:"pagination"`
}

func (s SearchCarriersRequest) Map() domain.CarrierSearchFilters {
	return domain.CarrierSearchFilters{
		Pagination: domain.Pagination{
			Page: s.Pagination.Page,
			Size: s.Pagination.PageSize,
		},
	}
}
