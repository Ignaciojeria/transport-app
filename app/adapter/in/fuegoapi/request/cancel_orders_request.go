package request

type CancelOrdersRequest struct {
	ManualChange struct {
		PerformedBy string `json:"performedBy" example:"juan@example.com"`
	} `json:"manualChange"`
	Orders []struct {
		BusinessIdentifiers struct {
			Commerce string `json:"commerce"`
			Consumer string `json:"consumer"`
		} `json:"businessIdentifiers"`
		ReferenceID string `json:"referenceID"`
	} `json:"orders"`
	CancellationReason struct {
		Detail      string `json:"detail"`
		Reason      string `json:"reason"`
		ReferenceID string `json:"referenceID"`
	} `json:"cancellationReason"`
}
