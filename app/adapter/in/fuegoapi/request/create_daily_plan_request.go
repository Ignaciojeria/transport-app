package request

import "transport-app/app/domain"

type CreateDailyPlanRequest struct {
	ReferenceID         string `json:"referenceID"`
	OperatorReferenceID string `json:"operatorReferenceID"`
	PlanDate            string `json:"planDate"`
	OrderReferenceIDs   []struct {
		ReferenceID string `json:"referenceID"`
	} `json:"orderReferenceIDs"`
}

func (r CreateDailyPlanRequest) Map() domain.Plan {
	var oders []domain.Order
	for _, v := range r.OrderReferenceIDs {
		oders = append(oders, domain.Order{
			ReferenceID: domain.ReferenceID(v.ReferenceID),
		})
	}
	return domain.Plan{
		ReferenceID: r.ReferenceID,
		Date:        r.PlanDate,
		PlanningStatus: domain.PlanningStatus{
			Value: "planned",
		},
		PlanType: domain.PlanType{
			Value: "dailyPlan",
		},
		Routes: []domain.Route{
			{
				Operator: domain.Operator{
					ReferenceID: r.OperatorReferenceID,
				},
				Orders: oders,
			},
		},
	}
}
