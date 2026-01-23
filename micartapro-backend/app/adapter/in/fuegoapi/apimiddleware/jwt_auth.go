package apimiddleware

import (
	"micartapro/app/shared/infrastructure/auth"
	"micartapro/app/shared/sharedcontext"
	"net/http"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type JWTAuthMiddleware func(http.Handler) http.Handler

func init() {
	ioc.Registry(NewJWTAuthMiddleware, auth.NewSupabaseTokenValidator)
}
func NewJWTAuthMiddleware(tokenValidator auth.SupabaseTokenValidator) JWTAuthMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")
			if token == "" {
				http.Error(w, "Bearer token is required", http.StatusUnauthorized)
				return
			}

			claims, err := tokenValidator.ValidateJWT(token)
			if err != nil {
				http.Error(w, "error validating token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Extraer el userId de las claims (en Supabase está en "sub")
			ctx := r.Context()
			if sub, ok := claims["sub"].(string); ok && sub != "" {
				ctx = sharedcontext.WithUserID(ctx, sub)
			}

			// Guardar el token en el contexto para uso posterior (ej. Storage API)
			ctx = sharedcontext.WithJWTToken(ctx, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
