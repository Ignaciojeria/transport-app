package tracking

import (
	"crypto/rand"
	"strings"
)

const crockfordBase32 = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"

// GenerateTrackingID genera un token opaco corto (8 chars) en Crockford Base32.
func GenerateTrackingID() (string, error) {
	const length = 8
	const entropyBytes = 5 // 5 bytes = 40 bits

	b := make([]byte, entropyBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		idx := int(b[i%entropyBytes]) & 31 // 0–31
		sb.WriteByte(crockfordBase32[idx])
	}

	return sb.String(), nil
}
