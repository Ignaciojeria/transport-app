package domain

import (
	"context"
	"fmt"
	"time"
)

type Webhook struct {
	Type      string            `json:"type"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      interface{}       `json:"body"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
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
	return nil
}

// UpdateIfChanged reemplaza el webhook por uno nuevo
func (w Webhook) UpdateIfChanged(newWebhook Webhook) (Webhook, bool) {
	return newWebhook, true
}
