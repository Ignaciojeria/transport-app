package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// ClientCredentialsEncryption maneja la encriptación/desencriptación de client secrets
type ClientCredentialsEncryption struct {
	key []byte
}

// NewClientCredentialsEncryption crea una nueva instancia del servicio de encriptación
func NewClientCredentialsEncryption(encryptionKey string) (*ClientCredentialsEncryption, error) {
	// Decodificar la clave base64
	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error decodificando clave de encriptación: %w", err)
	}

	// Verificar que la clave tenga el tamaño correcto para AES-256 (32 bytes)
	if len(key) != 32 {
		return nil, fmt.Errorf("la clave debe tener 32 bytes para AES-256")
	}

	return &ClientCredentialsEncryption{
		key: key,
	}, nil
}

// Encrypt encripta un texto usando AES-256-GCM
func (e *ClientCredentialsEncryption) Encrypt(plaintext string) (string, error) {
	// Crear cipher block
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("error creando cipher: %w", err)
	}

	// Crear GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creando GCM: %w", err)
	}

	// Crear nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generando nonce: %w", err)
	}

	// Encriptar
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Codificar en base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt desencripta un texto usando AES-256-GCM
func (e *ClientCredentialsEncryption) Decrypt(encryptedText string) (string, error) {
	// Decodificar de base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("error decodificando texto encriptado: %w", err)
	}

	// Crear cipher block
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("error creando cipher: %w", err)
	}

	// Crear GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creando GCM: %w", err)
	}

	// Verificar tamaño mínimo
	if len(ciphertext) < gcm.NonceSize() {
		return "", fmt.Errorf("texto encriptado demasiado corto")
	}

	// Extraer nonce
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Desencriptar
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error desencriptando: %w", err)
	}

	return string(plaintext), nil
}
