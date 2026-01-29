package fuegoapi

import (
	"encoding/json"
	"fmt"
	"os"
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type GenerateUploadURLRequest struct {
	FileName   string `json:"fileName"`   // Nombre del archivo (ej: "photo.jpg")
	ContentType string `json:"contentType"` // Tipo MIME (ej: "image/jpeg")
}

type GenerateUploadURLResponse struct {
	UploadURL  string `json:"uploadUrl"`  // URL firmada para subir
	PublicURL  string `json:"publicUrl"`   // URL pública de la imagen (bucket público)
	ObjectPath string `json:"objectPath"`  // Ruta del objeto en GCS
}

func init() {
	ioc.Registry(
		generateImageUploadURL,
		httpserver.New,
		gcs.NewClient,
		observability.NewObservability,
		configuration.NewConf,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func generateImageUploadURL(
	s httpserver.Server,
	gcsClient *storage.Client,
	obs observability.Observability,
	conf configuration.Conf,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Post(s.Manager, "/api/images/upload-url",
		func(c fuego.ContextWithBody[GenerateUploadURLRequest]) (GenerateUploadURLResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "generateImageUploadURL")
			defer span.End()

			// Extraer userID del contexto
			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return GenerateUploadURLResponse{}, fuego.HTTPError{
					Title:  "user_id not found in context",
					Detail: "user_id is required for image upload",
					Status: http.StatusUnauthorized,
				}
			}

			// Obtener el body del request
			req, err := c.Body()
			if err != nil {
				return GenerateUploadURLResponse{}, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Validar que fileName y contentType estén presentes
			if req.FileName == "" {
				return GenerateUploadURLResponse{}, fuego.HTTPError{
					Title:  "fileName is required",
					Detail: "fileName cannot be empty",
					Status: http.StatusBadRequest,
				}
			}

			if req.ContentType == "" {
				req.ContentType = "image/jpeg" // Default
			}

			// Construir la ruta del objeto: userID/timestamp-random-filename
			timestamp := time.Now().Unix()
			randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
			objectPath := fmt.Sprintf("%s/%d-%s-%s", userID, timestamp, randomSuffix, req.FileName)

			// Bucket de imágenes (puedes cambiarlo según tu configuración)
			bucketName := "micartapro-images" // O usa una variable de entorno
			
			// Generar signed URL para PUT (subir) usando storage.SignedURL
			// La URL será válida por 15 minutos
			// Necesitamos obtener las credenciales del service account
			opts := &storage.SignedURLOptions{
				Method:      "PUT",
				Expires:     time.Now().Add(15 * time.Minute),
				ContentType: req.ContentType,
			}

			// Intentar leer las credenciales desde GOOGLE_APPLICATION_CREDENTIALS
			credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
			if credsPath != "" {
				credsData, err := os.ReadFile(credsPath)
				if err == nil {
					var credsJSON map[string]interface{}
					if err := json.Unmarshal(credsData, &credsJSON); err == nil {
						if email, ok := credsJSON["client_email"].(string); ok {
							opts.GoogleAccessID = email
						}
						if key, ok := credsJSON["private_key"].(string); ok {
							opts.PrivateKey = []byte(key)
						}
					}
				}
			}

			// Si no se pudieron obtener las credenciales, retornar error
			if opts.GoogleAccessID == "" || opts.PrivateKey == nil {
				obs.Logger.ErrorContext(spanCtx, "missing_gcs_credentials", "error", "GoogleAccessID or PrivateKey not found")
				return GenerateUploadURLResponse{}, fuego.HTTPError{
					Title:  "GCS credentials not configured",
					Detail: "GOOGLE_APPLICATION_CREDENTIALS must point to a service account JSON file with client_email and private_key",
					Status: http.StatusInternalServerError,
				}
			}

			uploadURL, err := storage.SignedURL(bucketName, objectPath, opts)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_generating_signed_url", "error", err, "objectPath", objectPath)
				return GenerateUploadURLResponse{}, fuego.HTTPError{
					Title:  "error generating signed URL",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			// Construir la URL pública directa (bucket público)
			// Formato: https://storage.googleapis.com/<bucket>/<object-path>
			publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)

			obs.Logger.InfoContext(spanCtx, "signed_url_generated", 
				"objectPath", objectPath, 
				"userID", userID,
				"contentType", req.ContentType)

			return GenerateUploadURLResponse{
				UploadURL:  uploadURL,
				PublicURL:  publicURL,
				ObjectPath: objectPath,
			}, nil
		},
		option.Summary("Generate signed URL for image upload"),
		option.Description("Generates a signed URL that allows direct upload to Google Cloud Storage from the browser"),
		option.Tags("images"),
		option.Middleware(jwtAuthMiddleware),
	)
}
