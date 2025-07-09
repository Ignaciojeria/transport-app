package encryption

import (
	"testing"
)

func TestClientCredentialsEncryption(t *testing.T) {
	// Clave de prueba (32 bytes en base64)
	testKey := "dGVzdC1rZXktZm9yLWVuY3J5cHRpb24tdGVzdC1rZXk="

	// Crear servicio de encriptación
	encryptionService, err := NewClientCredentialsEncryption(testKey)
	if err != nil {
		t.Fatalf("Error creando servicio de encriptación: %v", err)
	}

	// Texto de prueba
	originalSecret := "mi-super-secreto-cliente-12345"

	// Encriptar
	encrypted, err := encryptionService.Encrypt(originalSecret)
	if err != nil {
		t.Fatalf("Error encriptando: %v", err)
	}

	// Verificar que el texto encriptado es diferente al original
	if encrypted == originalSecret {
		t.Error("El texto encriptado no debería ser igual al original")
	}

	// Desencriptar
	decrypted, err := encryptionService.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Error desencriptando: %v", err)
	}

	// Verificar que el texto desencriptado es igual al original
	if decrypted != originalSecret {
		t.Errorf("Texto desencriptado no coincide con el original. Esperado: %s, Obtenido: %s", originalSecret, decrypted)
	}
}

func TestClientCredentialsEncryptionWithDifferentInputs(t *testing.T) {
	// Clave de prueba
	testKey := "dGVzdC1rZXktZm9yLWVuY3J5cHRpb24tdGVzdC1rZXk="

	// Crear servicio de encriptación
	encryptionService, err := NewClientCredentialsEncryption(testKey)
	if err != nil {
		t.Fatalf("Error creando servicio de encriptación: %v", err)
	}

	// Casos de prueba
	testCases := []string{
		"",
		"a",
		"secret123",
		"mi-super-secreto-muy-largo-que-debe-funcionar-correctamente",
		"!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"áéíóúñ",
		"🔐🔒🔑",
	}

	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			// Encriptar
			encrypted, err := encryptionService.Encrypt(testCase)
			if err != nil {
				t.Fatalf("Error encriptando '%s': %v", testCase, err)
			}

			// Verificar que el texto encriptado es diferente al original
			if encrypted == testCase {
				t.Errorf("El texto encriptado no debería ser igual al original para '%s'", testCase)
			}

			// Desencriptar
			decrypted, err := encryptionService.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Error desencriptando '%s': %v", testCase, err)
			}

			// Verificar que el texto desencriptado es igual al original
			if decrypted != testCase {
				t.Errorf("Texto desencriptado no coincide con el original. Esperado: '%s', Obtenido: '%s'", testCase, decrypted)
			}
		})
	}
}

func TestClientCredentialsEncryptionInvalidKey(t *testing.T) {
	// Clave inválida (no es base64 válido)
	invalidKey := "invalid-key-not-base64"

	_, err := NewClientCredentialsEncryption(invalidKey)
	if err == nil {
		t.Error("Debería haber un error con una clave inválida")
	}
}

func TestClientCredentialsEncryptionWrongKeySize(t *testing.T) {
	// Clave con tamaño incorrecto (16 bytes en lugar de 32)
	shortKey := "dGVzdC1rZXk=" // 16 bytes en base64

	_, err := NewClientCredentialsEncryption(shortKey)
	if err == nil {
		t.Error("Debería haber un error con una clave de tamaño incorrecto")
	}
}
