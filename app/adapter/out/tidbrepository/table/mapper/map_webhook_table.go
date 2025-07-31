package mapper

import (
	"context"
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapWebhookToTable(ctx context.Context, w domain.Webhook) table.Webhook {
	// Serializar todo el webhook como JSON
	payloadBytes, _ := json.Marshal(w)
	payload := make(map[string]string)
	json.Unmarshal(payloadBytes, &payload)
	
	return table.Webhook{
		Payload:    table.JSONMap(payload),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		DocumentID: string(w.DocID(ctx)),
	}
} 