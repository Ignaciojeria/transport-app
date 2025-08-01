package sharedcontext

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/baggage"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	BaggageTenantID      = "tenant.id"
	BaggageTenantCountry = "tenant.country"
	BaggageAccountEmail  = "account.email"
	BaggageConsumer      = "consumer"
	BaggageCommerce      = "commerce"
	BaggageChannel       = "channel"
)

const (
	BaggageEventType  = "event.type"
	BaggageEntityType = "entity.type"
)

type EventContext struct {
	EventType  string
	EntityType string
	Consumer   string
	Commerce   string
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

func TenantIDFromContext(ctx context.Context) uuid.UUID {
	bag := baggage.FromContext(ctx)
	raw := bag.Member(BaggageTenantID).Value()
	if raw == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func TenantCountryFromContext(ctx context.Context) string {
	bag := baggage.FromContext(ctx)
	raw := bag.Member(BaggageTenantCountry).Value()
	return strings.ToUpper(raw)
}

// ChannelFromContext obtiene el channel desde el contexto
func ChannelFromContext(ctx context.Context) string {
	bag := baggage.FromContext(ctx)
	return bag.Member(BaggageChannel).Value()
}

// GetEventTypeFromContext obtiene el eventType desde el contexto
func GetEventTypeFromContext(ctx context.Context) (string, bool) {
	bag := baggage.FromContext(ctx)
	eventType := bag.Member(BaggageEventType).Value()
	return eventType, eventType != ""
}

// EntityTypeFromContext obtiene el entityType desde el contexto
func EntityTypeFromContext(ctx context.Context) string {
	bag := baggage.FromContext(ctx)
	return bag.Member(BaggageEntityType).Value()
}

const (
	accessTokenKey    string = "access.token"
	idempotencyKeyKey string = "idempotency.key"
	bucketTokenKey    string = "bucket.token"
)

// WithAccessToken añade el token de acceso al contexto de forma segura
func WithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, accessTokenKey, token)
}

// WithBucketToken añade el token del bucket al contexto de forma segura
func WithBucketToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, bucketTokenKey, token)
}

func WithIdempotencyKey(ctx context.Context, idempotencyKey string) context.Context {
	return context.WithValue(ctx, idempotencyKeyKey, idempotencyKey)
}

// AccessTokenFromContext obtiene el token de acceso desde el contexto
func AccessTokenFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(accessTokenKey)
	token, ok := val.(string)
	return token, ok
}

// BucketTokenFromContext obtiene el token del bucket desde el contexto
func BucketTokenFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(bucketTokenKey)
	token, ok := val.(string)
	return token, ok
}

func IdempotencyKeyFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(idempotencyKeyKey)
	key, ok := val.(string)
	return key, ok
}
