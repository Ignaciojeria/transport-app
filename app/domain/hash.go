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
func Hash(org Organization, inputs ...string) DocumentID {
	orgKey := org.GetOrgKey()
	joined := strings.Join(append([]string{orgKey}, inputs...), "|")
	hash := sha256.Sum256([]byte(joined))
	return DocumentID(hex.EncodeToString(hash[:16]))
}

// HashCtx construye un hash único usando tenantID, country y los inputs adicionales
func HashCtx(ctx context.Context, inputs ...string) DocumentID {
	bag := baggage.FromContext(ctx)

	tenantID := bag.Member(sharedcontext.BaggageTenantID).Value()
	country := bag.Member(sharedcontext.BaggageTenantCountry).Value()

	orgKey := tenantID + "-" + strings.ToUpper(country)

	joined := strings.Join(append([]string{orgKey}, inputs...), "|")
	hash := sha256.Sum256([]byte(joined))

	return DocumentID(hex.EncodeToString(hash[:16])) // truncate si lo necesitás corto
}
