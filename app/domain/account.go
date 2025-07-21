package domain

import (
	"crypto/sha256"

	"github.com/google/uuid"
)

type Account struct {
	Email string
	Role  string
}

func (a Account) DocID() DocumentID {
	return HashInputs(a.Email)
}

// UUID genera un UUID determinístico basado en el email del account
func (a Account) UUID() uuid.UUID {
	// Crear un namespace determinístico usando SHA-256 del string "transport-app-account"
	namespaceBytes := sha256.Sum256([]byte("transport-app-account"))

	// Crear un UUID válido para usar como namespace
	// Usar los primeros 16 bytes del hash SHA-256
	var namespace uuid.UUID
	copy(namespace[:], namespaceBytes[:16])

	// Asegurar que sea un UUID válido (versión 4, variante RFC4122)
	namespace[6] = (namespace[6] & 0x0f) | 0x40 // Version 4
	namespace[8] = (namespace[8] & 0x3f) | 0x80 // Variant RFC4122

	// Generar UUID v5 usando el namespace y el email
	return uuid.NewSHA1(namespace, []byte(a.Email))
}
