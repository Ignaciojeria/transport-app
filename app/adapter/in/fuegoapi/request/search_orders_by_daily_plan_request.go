package request

import (
	"time"
	"transport-app/app/domain"
)

type SearchOrdersByDailyPlanRequest struct {
	OperatorReferenceID string `json:"operatorReferenceID"`
	PlannedDate         string `json:"plannedDate"`
}

func (s SearchOrdersByDailyPlanRequest) Map() domain.OrderSearchFilters {
	plannedDate, err := time.Parse("2006-01-02", s.PlannedDate)
	if err != nil {
		plannedDate = time.Time{}
	}
	return domain.OrderSearchFilters{
		OrderSearchOperatorDailyPlanFilters: domain.OrderSearchOperatorDailyPlanFilters{
			OperatorReferenceID: s.OperatorReferenceID,
			PlannedDate:         plannedDate,
		},
	}
}
