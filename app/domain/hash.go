package domain

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"transport-app/app/shared/sharedcontext"

	"go.opentelemetry.io/otel/baggage"
)

// Hash genera un hash SHA-256 truncado a 128 bits (32 caracteres hex) y lo inicia con la key de la organización.
func HashByTenant(ctx context.Context, inputs ...string) DocumentID {
	bag := baggage.FromContext(ctx)

	tenantID := bag.Member(sharedcontext.BaggageTenantID).Value()
	country := bag.Member(sharedcontext.BaggageTenantCountry).Value()

	orgKey := tenantID + "-" + country

	joined := strings.Join(append([]string{orgKey}, inputs...), "|")
	hash := sha256.Sum256([]byte(joined))

	return DocumentID(hex.EncodeToString(hash[:16])) // truncate si lo necesitás corto
}

func HashByCountry(ctx context.Context, inputs ...string) DocumentID {
	bag := baggage.FromContext(ctx)

	country := bag.Member(sharedcontext.BaggageTenantCountry).Value()

	joined := strings.Join(append([]string{country}, inputs...), "|")
	hash := sha256.Sum256([]byte(joined))

	return DocumentID(hex.EncodeToString(hash[:16])) // truncate si lo necesitás corto
}
