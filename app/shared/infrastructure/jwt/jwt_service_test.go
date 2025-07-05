package jwt

import (
	"testing"
	"time"
)

func TestJWTService_GenerateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", "test-issuer")

	scopes := []string{"read:orders", "write:orders"}
	context := map[string]string{
		"tenant_id": "12345",
		"user_type": "admin",
	}

	token, expiresAt, err := jwtService.GenerateToken("user123", scopes, context, "test-tenant", "https://api.test.com", 60)
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	if token == "" {
		t.Error("Token no debería estar vacío")
	}

	if expiresAt <= time.Now().Unix() {
		t.Error("ExpiresAt debería ser en el futuro")
	}

	// Verificar que el token se puede validar
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Error validando token: %v", err)
	}

	if claims.Sub != "user123" {
		t.Errorf("Subject esperado: user123, obtenido: %s", claims.Sub)
	}

	if len(claims.Scopes) != 2 {
		t.Errorf("Scopes esperados: 2, obtenidos: %d", len(claims.Scopes))
	}

	if claims.Context["tenant_id"] != "12345" {
		t.Errorf("Context tenant_id esperado: 12345, obtenido: %s", claims.Context["tenant_id"])
	}

	if claims.Issuer != "test-issuer" {
		t.Errorf("Issuer esperado: test-issuer, obtenido: %s", claims.Issuer)
	}
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", "test-issuer")

	// Token inválido
	_, err := jwtService.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Debería haber error con token inválido")
	}
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", "test-issuer")

	// Generar token con expiración muy corta
	token, _, err := jwtService.GenerateToken("user123", []string{}, map[string]string{}, "test-tenant", "https://api.test.com", 0)
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	// Esperar un poco para que expire
	time.Sleep(1 * time.Second)

	_, err = jwtService.ValidateToken(token)
	if err == nil {
		t.Error("Debería haber error con token expirado")
	}
}

func TestJWTService_RefreshToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", "test-issuer")

	// Generar token original
	originalToken, originalExpiresAt, err := jwtService.GenerateToken("user123", []string{"admin"}, map[string]string{"test": "value"}, "test-tenant", "https://api.test.com", 60)
	if err != nil {
		t.Fatalf("Error generando token original: %v", err)
	}

	// Refrescar token
	newToken, newExpiresAt, err := jwtService.RefreshToken(originalToken, 120)
	if err != nil {
		t.Fatalf("Error refrescando token: %v", err)
	}

	if newToken == originalToken {
		t.Error("El nuevo token debería ser diferente al original")
	}

	if newExpiresAt <= originalExpiresAt {
		t.Error("La nueva expiración debería ser mayor que la original")
	}

	// Verificar que el nuevo token es válido
	claims, err := jwtService.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("Error validando token refrescado: %v", err)
	}

	if claims.Sub != "user123" {
		t.Errorf("Subject esperado: user123, obtenido: %s", claims.Sub)
	}

	if claims.Scopes[0] != "admin" {
		t.Errorf("Scope esperado: admin, obtenido: %s", claims.Scopes[0])
	}

	if claims.Context["test"] != "value" {
		t.Errorf("Context esperado: value, obtenido: %s", claims.Context["test"])
	}
}

func TestGenerateSecretKey(t *testing.T) {
	secretKey, err := GenerateSecretKey(32)
	if err != nil {
		t.Fatalf("Error generando clave secreta: %v", err)
	}

	if len(secretKey) == 0 {
		t.Error("La clave secreta no debería estar vacía")
	}

	// Generar otra clave para verificar que son diferentes
	secretKey2, err := GenerateSecretKey(32)
	if err != nil {
		t.Fatalf("Error generando segunda clave secreta: %v", err)
	}

	if secretKey == secretKey2 {
		t.Error("Las claves secretas deberían ser diferentes")
	}
}

func TestJWTService_DifferentSecretKeys(t *testing.T) {
	jwtService1 := NewJWTService("secret-key-1", "test-issuer")
	jwtService2 := NewJWTService("secret-key-2", "test-issuer")

	// Generar token con el primer servicio
	token, _, err := jwtService1.GenerateToken("user123", []string{}, map[string]string{}, "test-tenant", "https://api.test.com", 60)
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	// Intentar validar con el segundo servicio (debería fallar)
	_, err = jwtService2.ValidateToken(token)
	if err == nil {
		t.Error("Debería haber error al validar token con clave secreta diferente")
	}
}

func TestJWTService_AudienceInToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key", "test-issuer")

	scopes := []string{"read:orders"}
	context := map[string]string{"test": "value"}
	expectedAudience := "https://api.transport-app.com"

	token, _, err := jwtService.GenerateToken("user123", scopes, context, "test-tenant", expectedAudience, 60)
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	// Validar el token y verificar que el audience está presente
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Error validando token: %v", err)
	}

	if len(claims.Audience) == 0 {
		t.Error("El token debería tener un audience")
	}

	if claims.Audience[0] != expectedAudience {
		t.Errorf("Audience esperado: %s, obtenido: %s", expectedAudience, claims.Audience[0])
	}

	// Verificar que el audience está en el token decodificado
	if len(claims.RegisteredClaims.Audience) == 0 {
		t.Error("El token debería tener un audience en RegisteredClaims")
	}

	if claims.RegisteredClaims.Audience[0] != expectedAudience {
		t.Errorf("Audience en RegisteredClaims esperado: %s, obtenido: %s", expectedAudience, claims.RegisteredClaims.Audience[0])
	}
}
