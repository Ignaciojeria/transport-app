package jwt

import (
	"context"
	"net/http"
	"strings"
)

// JWTMiddleware crea un middleware para validar tokens JWT
func JWTMiddleware(jwtService *JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener el token del header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header requerido", http.StatusUnauthorized)
				return
			}

			// Verificar formato "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Formato de Authorization header inválido. Use: Bearer <token>", http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]

			// Validar token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Agregar claims al contexto de la request para uso posterior
			ctx := r.Context()
			ctx = context.WithValue(ctx, "jwt_claims", claims)
			ctx = context.WithValue(ctx, "jwt_sub", claims.Sub)
			ctx = context.WithValue(ctx, "jwt_scopes", claims.Scopes)
			ctx = context.WithValue(ctx, "jwt_context", claims.Context)

			// Crear nueva request con el contexto actualizado
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireScope crea un middleware que requiere un scope específico
func RequireScope(requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scopes, ok := r.Context().Value("jwt_scopes").([]string)
			if !ok {
				http.Error(w, "No se encontraron scopes en el token", http.StatusForbidden)
				return
			}

			// Verificar si el scope requerido está presente
			hasScope := false
			for _, scope := range scopes {
				if scope == requiredScope {
					hasScope = true
					break
				}
			}

			if !hasScope {
				http.Error(w, "Scope requerido no encontrado: "+requiredScope, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClaimsFromContext obtiene las claims JWT del contexto
func GetClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value("jwt_claims").(*Claims)
	return claims, ok
}

// GetSubjectFromContext obtiene el subject del token JWT del contexto
func GetSubjectFromContext(ctx context.Context) (string, bool) {
	sub, ok := ctx.Value("jwt_sub").(string)
	return sub, ok
}

// GetScopesFromContext obtiene los scopes del token JWT del contexto
func GetScopesFromContext(ctx context.Context) ([]string, bool) {
	scopes, ok := ctx.Value("jwt_scopes").([]string)
	return scopes, ok
}
