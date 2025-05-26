package projectionresult

type DeliveryUnitsProjectionResult struct {
	ID               int64  `json:"id"`
	OrderReferenceID string `json:"order_reference_id"`
	Channel          string `json:"channel"`
}
