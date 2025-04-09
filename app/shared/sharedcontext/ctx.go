// sharedcontext/baggage_keys.go
package sharedcontext

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/baggage"
)

const (
	BaggageTenantID      = "tenant.id"
	BaggageTenantCountry = "tenant.country"
	BaggageConsumer      = "business.consumer"
	BaggageCommerce      = "business.commerce"
)

const (
	BaggageEventType  = "event.type"
	BaggageEntityType = "entity.type"
)

type EventContext struct {
	EventType  string
	EntityType string
}

// AddEventContextToBaggage añade eventType y entityType al baggage del contexto
func AddEventContextToBaggage(ctx context.Context, evt EventContext) context.Context {
	m1, _ := baggage.NewMember(BaggageEventType, string(evt.EventType))
	m2, _ := baggage.NewMember(BaggageEntityType, string(evt.EntityType))

	existing := baggage.FromContext(ctx)
	allMembers := append(existing.Members(), m1, m2)

	bag, _ := baggage.New(allMembers...)
	return baggage.ContextWithBaggage(ctx, bag)
}

// CopyBaggageToAttributesCamelCase transforma keys tipo "tenant.id" -> "tenantId"
func CopyBaggageToAttributesCamelCase(ctx context.Context, attrs map[string]string) {
	bag := baggage.FromContext(ctx)

	for _, m := range bag.Members() {
		// Convierte la clave a camelCase si tiene punto
		keyParts := strings.Split(m.Key(), ".")
		if len(keyParts) == 1 {
			attrs[m.Key()] = m.Value()
			continue
		}

		// Ejemplo: ["tenant", "id"] => "tenantId"
		camelKey := keyParts[0]
		for _, part := range keyParts[1:] {
			if len(part) > 0 {
				camelKey += strings.Title(part) // o usar capitalizeFirst(part)
			}
		}

		attrs[camelKey] = m.Value()
	}
}
