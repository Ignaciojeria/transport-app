package sharedcontext

import (
	"context"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/baggage"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	BaggageTenantID      = "tenant.id"
	BaggageTenantCountry = "tenant.country"
	BaggageAccountEmail  = "account.email"
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
	capitalizer := cases.Title(language.English)

	for _, m := range bag.Members() {
		keyParts := strings.Split(m.Key(), ".")
		if len(keyParts) == 1 {
			attrs[m.Key()] = m.Value()
			continue
		}

		camelKey := keyParts[0]
		for _, part := range keyParts[1:] {
			if len(part) > 0 {
				camelKey += capitalizer.String(part)
			}
		}

		attrs[camelKey] = m.Value()
	}
}

func TenantIDFromContext(ctx context.Context) int64 {
	bag := baggage.FromContext(ctx)
	raw := bag.Member(BaggageTenantID).Value()
	if raw == "" {
		return 0
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func TenantCountryFromContext(ctx context.Context) string {
	bag := baggage.FromContext(ctx)
	raw := bag.Member(BaggageTenantCountry).Value()
	return strings.ToUpper(raw)
}
