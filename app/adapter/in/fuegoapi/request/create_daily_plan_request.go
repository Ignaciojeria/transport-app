package request

type CreateDailyPlanRequest struct {
	ReferenceID       string `json:"referenceID"`
	PlanDate          string `json:"planDate"`
	OrderReferenceIDs []struct {
		ReferenceID string `json:"referenceID"`
	} `json:"orderReferenceIDs"`
}
