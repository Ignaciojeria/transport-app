package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Webhook struct {
	Type        string            `json:"type"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	RetryPolicy WebhookRetryPolicy `json:"retryPolicy"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type WebhookRetryPolicy struct {
	MaxRetries     int `json:"maxRetries"`
	BackoffSeconds int `json:"backoffSeconds"`
}

func (w Webhook) DocID(ctx context.Context) DocumentID {
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%s", w.Type, w.URL, traceID)))
	return DocumentID{
		Value: hex.EncodeToString(hash[:]),
	}
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