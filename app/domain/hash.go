package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Hash genera un hash SHA-256 truncado a 128 bits (32 caracteres hex) y lo inicia con la key de la organización.
func Hash(org Organization, inputs ...string) string {
	// Obtener la clave de la organización
	orgKey := org.GetOrgKey()

	// Unir todos los inputs asegurando que la key de la organización esté al inicio
	joined := strings.Join(append([]string{orgKey}, inputs...), "|")

	// Generar el hash SHA-256
	hash := sha256.Sum256([]byte(joined))

	// Retornar los primeros 128 bits (32 caracteres hex)
	return hex.EncodeToString(hash[:16])
}
