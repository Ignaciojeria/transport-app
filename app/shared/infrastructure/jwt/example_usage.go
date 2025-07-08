package jwt

import (
	"fmt"
	"log"
)

// ExampleUsage muestra cómo usar el sistema JWT completo
func exampleUsage() {
	// 1. Crear el servicio JWT
	jwtService := NewJWTService("mi-clave-secreta-super-segura", "transport-app")

	// 2. Generar un token
	scopes := []string{"read:orders", "write:orders", "admin"}
	context := map[string]string{
		"tenant_id": "12345",
		"user_type": "admin",
	}

	token, err := jwtService.GenerateToken(
		"user123",
		scopes,
		context,
		"tenant123",                     // tenant ID
		"https://api.transport-app.com", // audience
		60,                              // 60 minutos
	)
	if err != nil {
		log.Fatal("Error generando token:", err)
	}

	fmt.Printf("Token generado: %s\n", token)

	// 3. Validar el token
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		log.Fatal("Error validando token:", err)
	}

	fmt.Printf("Token válido para usuario: %s\n", claims.Sub)
	fmt.Printf("Scopes: %v\n", claims.Scopes)
	fmt.Printf("Context: %v\n", claims.Context)
	fmt.Printf("Expira en: %d\n", claims.ExpiresAt.Unix())

	// 4. Refrescar el token
	newToken, err := jwtService.RefreshToken(token, 120) // 2 horas
	if err != nil {
		log.Fatal("Error refrescando token:", err)
	}

	fmt.Printf("Nuevo token: %s\n", newToken)

	// Validar el nuevo token para obtener su expiración
	newClaims, err := jwtService.ValidateToken(newToken)
	if err != nil {
		log.Fatal("Error validando nuevo token:", err)
	}
	fmt.Printf("Nueva expiración: %d\n", newClaims.ExpiresAt.Unix())

	// 5. Generar una clave secreta aleatoria
	secretKey, err := GenerateSecretKey(32)
	if err != nil {
		log.Fatal("Error generando clave secreta:", err)
	}

	fmt.Printf("Clave secreta generada: %s\n", secretKey)
}

// ExampleMiddlewareUsage muestra cómo usar el middleware
func ExampleMiddlewareUsage() {
	// Este ejemplo muestra cómo se usaría en un servidor HTTP real
	// con el middleware JWT

	jwtService := NewJWTService("mi-clave-secreta", "transport-app")

	// Crear middleware JWT
	_ = JWTMiddleware(jwtService)

	// Crear middleware para requerir scope específico
	_ = RequireScope("admin")

	// En un servidor real, se usarían así:
	// server.Use(jwtMiddleware)
	// server.Use(adminScopeMiddleware)

	fmt.Println("Middleware JWT configurado correctamente")
	fmt.Println("Para usar en rutas protegidas:")
	fmt.Println("1. Agregar JWTMiddleware para validar tokens")
	fmt.Println("2. Agregar RequireScope para verificar permisos específicos")
	fmt.Println("3. Usar GetClaimsFromContext para obtener información del token")
}
