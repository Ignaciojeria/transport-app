package request

type UpsertPlanRequest struct {
	StartDate           string   `json:"startDate"`
	EndDate             string   `json:"endDate"`
	VehiclePlates       []string `json:"vehiclePlates"`
	OrderReferenceIDs   []string `json:"orderReferenceIDs"`
	OriginNodeReference string   `json:"originNodeReferenceID"`
}
