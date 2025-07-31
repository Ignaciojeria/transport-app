package domain

import (
	"context"
	"fmt"
	"time"
)

type Webhook struct {
	Type        string             `json:"type"`
	URL         string             `json:"url"`
	Headers     map[string]string  `json:"headers"`
	RetryPolicy WebhookRetryPolicy `json:"retryPolicy"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type WebhookRetryPolicy struct {
	MaxRetries     int `json:"maxRetries"`
	BackoffSeconds int `json:"backoffSeconds"`
}

func (w Webhook) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, w.Type)
}

func (w Webhook) Validate() error {
	if w.Type == "" {
		return fmt.Errorf("webhook type is required")
	}
	if w.URL == "" {
		return fmt.Errorf("webhook URL is required")
	}
	if w.RetryPolicy.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}
	if w.RetryPolicy.BackoffSeconds < 0 {
		return fmt.Errorf("backoff seconds cannot be negative")
	}
	return nil
}

func (w Webhook) UpdateIfChanged(newWebhook Webhook) (Webhook, bool) {
	// Para webhooks, siempre actualizar con el nuevo valor
	return newWebhook, true
}
