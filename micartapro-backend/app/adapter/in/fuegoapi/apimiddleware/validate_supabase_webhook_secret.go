package apimiddleware

import (
	"micartapro/app/shared/configuration"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type ValidateSupabaseWebhookSecretMiddleware func(http.Handler) http.Handler

func init() {
	ioc.Registry(NewValidateSupabaseWebhookSecretMiddleware, configuration.NewConf)
}

func NewValidateSupabaseWebhookSecretMiddleware(conf configuration.Conf) ValidateSupabaseWebhookSecretMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := r.Header.Get("X-Supabase-Webhook-Secret")
			if secret != conf.SUPABASE_WEBHOOK_SECRET {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
