package apimiddleware

import (
	"micartapro/app/shared/infrastructure/auth"
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

			_, err := tokenValidator.ValidateJWT(token)
			if err != nil {
				http.Error(w, "error validating token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
