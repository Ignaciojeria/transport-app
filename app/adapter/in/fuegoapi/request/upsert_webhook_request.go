package request

type UpsertWebhookRequest struct {
	Type        string            `json:"type"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	RetryPolicy struct {
		MaxRetries     int `json:"maxRetries"`
		BackoffSeconds int `json:"backoffSeconds"`
	} `json:"retryPolicy"`
}
