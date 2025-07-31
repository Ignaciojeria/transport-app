package domain

import "github.com/google/uuid"

type Webhook struct {
	ID          uuid.UUID         `json:"id"`
	Type        string            `json:"type"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	RetryPolicy RetryPolicy       `json:"retryPolicy"`
}

type RetryPolicy struct {
	MaxRetries     int `json:"maxRetries"`
	BackoffSeconds int `json:"backoffSeconds"`
}