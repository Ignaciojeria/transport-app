package creem

import (
	"encoding/json"
	"time"

	"micartapro/app/usecase/billing"

	"github.com/google/uuid"
)

// extractSubscriptionID busca el subscription_id en diferentes lugares del objeto
func extractSubscriptionID(object json.RawMessage) string {
	var obj map[string]interface{}
	if err := json.Unmarshal(object, &obj); err != nil {
		return ""
	}

	// Si el objeto es la subscription directamente (tiene id y status)
	if id, ok := obj["id"].(string); ok {
		if status, ok := obj["status"].(string); ok && status != "" {
			return id
		}
	}

	// Si la subscription está anidada
	if subscription, ok := obj["subscription"].(map[string]interface{}); ok {
		if id, ok := subscription["id"].(string); ok {
			return id
		}
	}

	return ""
}

// toBillingEvent convierte un webhook de Creem a BillingEvent
func toBillingEvent(webhookID string, eventType string, createdAt int64, payload json.RawMessage, userID *uuid.UUID) billing.BillingEvent {
	// Convertir CreatedAt de milisegundos a time.Time
	providerCreatedAt := time.Unix(0, createdAt*int64(time.Millisecond))

	// Extraer subscription_id del payload
	var webhookData map[string]interface{}
	subscriptionID := ""
	if err := json.Unmarshal(payload, &webhookData); err == nil {
		if object, ok := webhookData["object"].(map[string]interface{}); ok {
			objectBytes, _ := json.Marshal(object)
			subscriptionID = extractSubscriptionID(json.RawMessage(objectBytes))
		}
	}

	return billing.BillingEvent{
		Provider:          "creem",
		ProviderEventID:   webhookID,
		EventType:         eventType,
		SubscriptionID:    subscriptionID,
		UserID:            userID,
		Payload:           payload,
		ProviderCreatedAt: providerCreatedAt,
	}
}
