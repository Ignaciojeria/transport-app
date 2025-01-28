package response

type SearchFirstOperatorDailyPlanResponse struct {
	ReferenceID     string `json:"referenceID"`
	PlanDate        string `json:"planDate"`
	OrderReferences []struct {
		ReferenceID          string  `json:"referenceID"`
		SequenceNumber       int     `json:"sequenceNumber"`
		DeliveryInstructions string  `json:"deliveryInstructions"`
		CustomerName         string  `json:"customerName"`
		CustomerAddress      string  `json:"customerAddress,omitempty"`
		PackageLPN           string  `json:"packageLPN,omitempty"`
		Latitude             float64 `json:"latitude"`
		Longitude            float64 `json:"longitude"`
		Items                []struct {
			Description string `json:"description"`
			Quantity    int    `json:"quantity"`
		} `json:"items"`
	} `json:"orderReferences"`
}
