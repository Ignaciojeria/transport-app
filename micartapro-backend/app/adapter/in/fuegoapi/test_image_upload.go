package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/imageuploader"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"
	"time"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Register(testImageUpload)
}

func testImageUpload(
	s httpserver.Server,
	uploadImage imageuploader.UploadImage,
	obs observability.Observability,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/test-image-upload",
		func(c fuego.ContextNoBody) (map[string]any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "testImageUpload")
			defer span.End()

			// Extraer userID del contexto (requerido por el uploader)
			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return nil, fuego.HTTPError{
					Title:  "user_id not found in context",
					Detail: "user_id is required for image upload",
					Status: http.StatusUnauthorized,
				}
			}

			// Crear un texto simple en duro para probar
			testText := "Este es un texto de prueba para subir a Supabase Storage. Fecha: " + time.Now().Format(time.RFC3339)
			testBytes := []byte(testText)

			// Generar un nombre de archivo único
			fileName := "test-upload-" + time.Now().Format("20060102-150405") + ".txt"

			obs.Logger.InfoContext(spanCtx, "testing_image_upload", "fileName", fileName, "size_bytes", len(testBytes), "userID", userID)

			// Subir el archivo usando el image uploader
			publicURL, err := uploadImage(spanCtx, testBytes, fileName)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_uploading_test_file", "error", err, "fileName", fileName)
				return nil, fuego.HTTPError{
					Title:  "error uploading file",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "test_file_uploaded_successfully", "publicURL", publicURL, "fileName", fileName)

			return map[string]any{
				"success":   true,
				"fileName":  fileName,
				"publicURL": publicURL,
				"sizeBytes": len(testBytes),
				"userID":    userID,
			}, nil
		},
		option.Summary("testImageUpload"),
		option.Tags("test"),
		option.Middleware(jwtAuthMiddleware),
	)
}
