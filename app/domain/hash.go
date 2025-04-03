package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Hash genera un hash SHA-256 truncado a 128 bits (32 caracteres hex) y lo inicia con la key de la organización.
func Hash(org Organization, inputs ...string) DocumentID {
	orgKey := org.GetOrgKey()
	joined := strings.Join(append([]string{orgKey}, inputs...), "|")
	hash := sha256.Sum256([]byte(joined))
	return DocumentID(hex.EncodeToString(hash[:16]))
}
