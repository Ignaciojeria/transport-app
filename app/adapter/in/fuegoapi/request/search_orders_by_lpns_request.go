package request

import "transport-app/app/domain"

type SearchOrdersByLpnsRequest struct {
	Commerces   []string `json:"commerces"`
	PackageLpns []string `json:"packageLpns"`
}

func (r SearchOrdersByLpnsRequest) Map() domain.OrderSearchFilters {
	return domain.OrderSearchFilters{
		Lpns:      r.PackageLpns,
		Commerces: r.Commerces,
	}
}
