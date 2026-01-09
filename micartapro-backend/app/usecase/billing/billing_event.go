package billing

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BillingEvent struct {
	// Provider
	Provider string // "creem"

	// Idempotencia fuerte
	ProviderEventID string // evt_...

	// Tipo semántico del evento
	EventType string // subscription.paid, checkout.completed, etc.

	// Ordering key / lifecycle
	SubscriptionID string // sub_...

	// Asociación a tu dominio (desde metadata.user_id)
	UserID *uuid.UUID

	// Payload completo e INMUTABLE
	Payload json.RawMessage

	// Timestamp del evento según Creem
	ProviderCreatedAt time.Time
}
