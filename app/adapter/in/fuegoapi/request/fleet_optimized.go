package request

type FleetOptimizedWebhookBody struct {
	Plan   string   `json:"plan"`
	Routes []string `json:"routes"`
}
