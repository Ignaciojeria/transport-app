package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Webhook struct {
	ID          uuid.UUID         `json:"id"`
	Type        string            `json:"type"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	RetryPolicy RetryPolicy       `json:"retryPolicy"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type RetryPolicy struct {
	MaxRetries     int `json:"maxRetries"`
	BackoffSeconds int `json:"backoffSeconds"`
}

// DocID genera un identificador único para el webhook basado en el contexto y tipo
func (w Webhook) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, w.Type)
}

// Validate valida los campos obligatorios del webhook
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

// UpdateIfChanged reemplaza el webhook por uno nuevo
func (w Webhook) UpdateIfChanged(newWebhook Webhook) (Webhook, bool) {
	return newWebhook, true
}