package request

import (
	"time"
	"transport-app/app/domain"
)

type UpsertDailyPlanRequest struct {
	OperatorReferenceID string `json:"operatorReferenceID"`
	PlannedDate         string `json:"plannedDate"`
	OrderReferenceIDs   []struct {
		ReferenceID string `json:"referenceID"`
	} `json:"orderReferenceIDs"`
}

func (r UpsertDailyPlanRequest) Map() domain.Plan {
	// Convertir string a time.Time
	planDate, err := time.Parse("2006-01-02", r.PlannedDate)
	if err != nil {
		// Dependiendo de tu manejo de errores podrías:
		// 1. Retornar un error adicional en la firma del método
		// 2. Usar un valor por defecto
		// 3. Usar time.Time{} (zero value)
		planDate = time.Time{} // usando zero value como ejemplo
	}

	var oders []domain.Order
	for _, v := range r.OrderReferenceIDs {
		oders = append(oders, domain.Order{
			ReferenceID: domain.ReferenceID(v.ReferenceID),
		})
	}

	return domain.Plan{
		ReferenceID: r.ReferenceID(),
		PlannedDate: planDate, // Usar el time.Time en lugar del string
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

func (r UpsertDailyPlanRequest) ReferenceID() string {
	return r.PlannedDate + "_" + r.OperatorReferenceID
}
