package apimiddleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
)

type ValidateCreemWebhookSecretMiddleware func(http.Handler) http.Handler

func init() {
	ioc.Register(NewValidateCreemWebhookSecretMiddleware)
}

func NewValidateCreemWebhookSecretMiddleware(conf configuration.Conf, obs observability.Observability) ValidateCreemWebhookSecretMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Leer el cuerpo de la solicitud
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
				return
			}
			// Restaurar el cuerpo de la solicitud para que pueda ser leído nuevamente
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Imprimir el body para debugging
			obs.Logger.Info("creem webhook body received",
				"body", string(body),
				"path", r.URL.Path,
			)

			// Obtener la firma del encabezado
			signature := r.Header.Get("Creem-Signature")
			if signature == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validar que el secreto esté configurado
			if conf.CREEM_WEBHOOK_SIGNING_SECRET == "" {
				http.Error(w, "Secreto de firma no configurado", http.StatusInternalServerError)
				return
			}

			// Calcular la firma HMAC-SHA256
			mac := hmac.New(sha256.New, []byte(conf.CREEM_WEBHOOK_SIGNING_SECRET))
			mac.Write(body)
			expectedMAC := mac.Sum(nil)

			// Comparar la firma calculada con la proporcionada (usar hmac.Equal para evitar timing attacks)
			signatureBytes, err := hex.DecodeString(signature)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !hmac.Equal(expectedMAC, signatureBytes) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
