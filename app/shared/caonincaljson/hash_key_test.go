package canonicaljson

import (
	"context"
	"testing"
	"transport-app/app/shared/sharedcontext"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/baggage"
)

type testPayload struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type testPayloadReversed struct {
	B string `json:"b"`
	A int    `json:"a"`
}

func TestHashKey_Deterministic(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()
	tID, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID.String())
	bag, _ := baggage.New(tID)
	ctx = baggage.ContextWithBaggage(ctx, bag)

	p1 := testPayload{A: 1, B: "hola"}
	p2 := testPayloadReversed{B: "hola", A: 1}

	hash1, err := HashKey(ctx, "prefix", p1)
	require.NoError(t, err)

	hash2, err := HashKey(ctx, "prefix", p2)
	require.NoError(t, err)

	require.Equal(t, hash1, hash2, "El hash debe ser igual para payloads equivalentes")
}

func TestHashKey_DifferentPayloads(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()
	tID, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID.String())
	bag, _ := baggage.New(tID)
	ctx = baggage.ContextWithBaggage(ctx, bag)

	p1 := testPayload{A: 1, B: "hola"}
	p2 := testPayload{A: 2, B: "hola"}

	hash1, err := HashKey(ctx, "prefix", p1)
	require.NoError(t, err)

	hash2, err := HashKey(ctx, "prefix", p2)
	require.NoError(t, err)

	require.NotEqual(t, hash1, hash2, "El hash debe ser diferente para payloads distintos")
}

func TestHashKey_DifferentTenants(t *testing.T) {
	p1 := testPayload{A: 1, B: "hola"}

	// Primer tenant
	ctx1 := context.Background()
	tenantID1 := uuid.New()
	tID1, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID1.String())
	bag1, _ := baggage.New(tID1)
	ctx1 = baggage.ContextWithBaggage(ctx1, bag1)

	hash1, err := HashKey(ctx1, "prefix", p1)
	require.NoError(t, err)

	// Segundo tenant
	ctx2 := context.Background()
	tenantID2 := uuid.New()
	tID2, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID2.String())
	bag2, _ := baggage.New(tID2)
	ctx2 = baggage.ContextWithBaggage(ctx2, bag2)

	hash2, err := HashKey(ctx2, "prefix", p1)
	require.NoError(t, err)

	require.NotEqual(t, hash1, hash2, "El hash debe ser diferente para tenants distintos")
}

func TestHashKey_SameTenant(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()
	tID, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID.String())
	bag, _ := baggage.New(tID)
	ctx = baggage.ContextWithBaggage(ctx, bag)

	p1 := testPayload{A: 1, B: "hola"}

	hash1, err := HashKey(ctx, "prefix", p1)
	require.NoError(t, err)

	hash2, err := HashKey(ctx, "prefix", p1)
	require.NoError(t, err)

	require.Equal(t, hash1, hash2, "El hash debe ser igual para el mismo tenant y payload")
}
