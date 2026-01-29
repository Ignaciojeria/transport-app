package imageuploader

import (
	"context"
	"fmt"
	"time"

	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/infrastructure/supabasecli"
	"micartapro/app/shared/sharedcontext"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	supabase "github.com/supabase-community/supabase-go"
)

type UploadImage func(ctx context.Context, imageBytes []byte, fileName string) (string, error)

func init() {
	ioc.Registry(NewImageUploader, observability.NewObservability, gcs.NewClient, supabasecli.NewSupabaseClient)
}

func NewImageUploader(obs observability.Observability, gcsClient *storage.Client, supabaseClient *supabase.Client) (UploadImage, error) {
	return func(ctx context.Context, imageBytes []byte, fileName string) (string, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "upload_image")
		defer span.End()

		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			return "", fmt.Errorf("userID is required but not found in context")
		}

		obs.Logger.InfoContext(spanCtx, "uploading_image", "fileName", fileName, "size_bytes", len(imageBytes), "userID", userID)

		// Construir la ruta del objeto: userID/timestamp-random-filename
		timestamp := time.Now().Unix()
		randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
		objectPath := fmt.Sprintf("%s/%d-%s-%s", userID, timestamp, randomSuffix, fileName)

		// Bucket de imágenes en GCS
		bucketName := "micartapro-images"
		bucket := gcsClient.Bucket(bucketName)
		object := bucket.Object(objectPath)

		// Crear writer para subir la imagen
		writer := object.NewWriter(spanCtx)
		writer.ContentType = "image/jpeg" // Default, se puede mejorar detectando el tipo real
		writer.CacheControl = "public, max-age=31536000" // Cache por 1 año

		// Escribir los bytes de la imagen
		if _, err := writer.Write(imageBytes); err != nil {
			writer.Close()
			obs.Logger.ErrorContext(spanCtx, "error_writing_image", "error", err, "objectPath", objectPath)
			return "", fmt.Errorf("error writing image to GCS: %w", err)
		}

		// Cerrar el writer para completar la subida
		if err := writer.Close(); err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_closing_writer", "error", err, "objectPath", objectPath)
			return "", fmt.Errorf("error closing GCS writer: %w", err)
		}

		// Construir la URL pública de la imagen (bucket público)
		// Formato: https://storage.googleapis.com/<bucket>/<object-path>
		publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)

		// Guardar la URL en la tabla catalog_images de Supabase
		record := map[string]interface{}{
			"id":        uuid.New().String(),
			"image_url": publicURL,
			"user_id":   userID,
		}

		_, _, err := supabaseClient.From("catalog_images").
			Insert(record, false, "", "", "").
			Execute()

		if err != nil {
			// Log el error pero no fallar la subida de imagen
			// La imagen ya está en GCS, solo falló el registro en la BD
			obs.Logger.WarnContext(spanCtx, "error_saving_to_catalog_images", "error", err, "publicURL", publicURL, "userID", userID)
		} else {
			obs.Logger.InfoContext(spanCtx, "image_saved_to_catalog_images", "publicURL", publicURL, "userID", userID)
		}

		obs.Logger.InfoContext(spanCtx, "image_uploaded_successfully", "publicURL", publicURL, "fileName", fileName, "objectPath", objectPath, "size_bytes", len(imageBytes))

		return publicURL, nil
	}, nil
}
