package canonicaljson

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"transport-app/app/shared/sharedcontext"

	"olympos.io/encoding/cjson"
)

// HashKey genera una key determinística para el payload dado incluyendo el tenant del contexto
func HashKey(ctx context.Context, prefix string, payload interface{}) (string, error) {
	// Obtener el tenant del contexto
	tenantID := sharedcontext.TenantIDFromContext(ctx)

	// Crear un objeto que incluya el tenant y el payload
	hashData := struct {
		TenantID string      `json:"tenant_id"`
		Payload  interface{} `json:"payload"`
	}{
		TenantID: tenantID.String(),
		Payload:  payload,
	}

	jsonBytes, err := cjson.Marshal(hashData)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(jsonBytes)
	hash := hex.EncodeToString(sum[:])

	return fmt.Sprintf("%s:%s", prefix, hash), nil
}
