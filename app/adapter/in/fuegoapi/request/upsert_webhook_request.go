package request

import (
	"context"
	"time"
	"transport-app/app/domain"
)

type UpsertWebhookRequest struct {
	Type    string            `json:"type"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func (req UpsertWebhookRequest) Map(ctx context.Context) domain.Webhook {
	now := time.Now()
	return domain.Webhook{
		Type:      req.Type,
		URL:       req.URL,
		Headers:   req.Headers,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
