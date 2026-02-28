package apimiddleware

import (
	"micartapro/app/shared/sharedcontext"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
)

type IdempotencyKeyMiddleware func(http.Handler) http.Handler

func init() {
	ioc.Register(NewIdempotencyKeyMiddleware)
}

func NewIdempotencyKeyMiddleware() IdempotencyKeyMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if idempotencyKey := r.Header.Get("Idempotency-Key"); idempotencyKey != "" {
				// Validar que el idempotency key sea un UUIDv7 válido
				parsedUUID, err := uuid.Parse(idempotencyKey)
				if err != nil {
					http.Error(w, "Idempotency-Key must be a valid UUIDv7", http.StatusBadRequest)
					return
				}
				// Verificar que sea UUIDv7 (versión 7)
				if parsedUUID.Version() != 7 {
					http.Error(w, "Idempotency-Key must be a UUIDv7", http.StatusBadRequest)
					return
				}
				ctx = sharedcontext.WithIdempotencyKey(ctx, idempotencyKey)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
