package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapWebhookToTable(ctx context.Context, w domain.Webhook) (table.Webhook, error) {
	// Serialize the webhook to JSON properly
	payloadBytes, err := json.Marshal(w)
	if err != nil {
		return table.Webhook{}, fmt.Errorf("failed to marshal webhook to JSON: %w", err)
	}
	
	// Create a proper JSON map that preserves data types
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		return table.Webhook{}, fmt.Errorf("failed to unmarshal webhook JSON: %w", err)
	}
	
	// Convert to JSONMap format for storage
	payload := make(map[string]string)
	for key, value := range payloadMap {
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return table.Webhook{}, fmt.Errorf("failed to marshal field %s: %w", key, err)
		}
		payload[key] = string(valueBytes)
	}
	
	return table.Webhook{
		Payload:    table.JSONMap(payload),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		DocumentID: string(w.DocID(ctx)),
	}, nil
} 