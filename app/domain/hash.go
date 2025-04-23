package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"transport-app/app/shared/sharedcontext"

	"go.opentelemetry.io/otel/baggage"
)

// HashByTenant genera un hash SHA-256 completo con contexto de organización
func HashByTenant(ctx context.Context, inputs ...string) DocumentID {
	bag := baggage.FromContext(ctx)

	tenantID := bag.Member(sharedcontext.BaggageTenantID).Value()
	country := bag.Member(sharedcontext.BaggageTenantCountry).Value()

	orgKey := tenantID + "-" + country
	joined := strings.Join(append([]string{orgKey}, inputs...), "|")

	hash := sha256.Sum256([]byte(joined))
	return DocumentID(hex.EncodeToString(hash[:])) // NO truncado
}

// HashByCountry genera un hash SHA-256 completo con contexto de país
func HashByCountry(ctx context.Context, inputs ...string) DocumentID {
	bag := baggage.FromContext(ctx)

	country := bag.Member(sharedcontext.BaggageTenantCountry).Value()
	joined := strings.Join(append([]string{country}, inputs...), "|")

	hash := sha256.Sum256([]byte(joined))
	return DocumentID(hex.EncodeToString(hash[:])) // NO truncado
}
