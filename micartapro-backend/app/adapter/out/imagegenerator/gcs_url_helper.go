package imagegenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"micartapro/app/shared/infrastructure/observability"
)

// GenerateSignedReadURL genera una signed URL de lectura para una imagen en GCS
// Si la URL ya es pública o no es de GCS, la retorna sin modificar
func GenerateSignedReadURL(ctx context.Context, obs observability.Observability, imageURL string) (string, error) {
	// Si la URL no es de GCS (storage.googleapis.com), retornarla sin modificar
	if !strings.Contains(imageURL, "storage.googleapis.com") {
		obs.Logger.InfoContext(ctx, "url_not_gcs", "url", imageURL)
		return imageURL, nil
	}

	// Extraer bucket y object path de la URL
	// Formato: https://storage.googleapis.com/<bucket>/<object-path>
	parts := strings.Split(imageURL, "storage.googleapis.com/")
	if len(parts) != 2 {
		obs.Logger.WarnContext(ctx, "invalid_gcs_url_format", "url", imageURL)
		return imageURL, nil // Retornar URL original si no podemos parsearla
	}

	pathParts := strings.SplitN(parts[1], "/", 2)
	if len(pathParts) != 2 {
		obs.Logger.WarnContext(ctx, "invalid_gcs_path_format", "url", imageURL)
		return imageURL, nil
	}

	bucketName := pathParts[0]
	objectPath := pathParts[1]

	// Generar signed URL para lectura (válida por 1 hora)
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(1 * time.Hour),
	}

	// Obtener credenciales del service account
	credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsPath == "" {
		obs.Logger.WarnContext(ctx, "no_gcs_credentials", "message", "GOOGLE_APPLICATION_CREDENTIALS not set, using original URL")
		return imageURL, nil
	}

	credsData, err := os.ReadFile(credsPath)
	if err != nil {
		obs.Logger.WarnContext(ctx, "error_reading_credentials", "error", err)
		return imageURL, nil
	}

	var credsJSON map[string]interface{}
	if err := json.Unmarshal(credsData, &credsJSON); err != nil {
		obs.Logger.WarnContext(ctx, "error_parsing_credentials", "error", err)
		return imageURL, nil
	}

	email, ok := credsJSON["client_email"].(string)
	if !ok || email == "" {
		obs.Logger.WarnContext(ctx, "missing_client_email")
		return imageURL, nil
	}

	key, ok := credsJSON["private_key"].(string)
	if !ok || key == "" {
		obs.Logger.WarnContext(ctx, "missing_private_key")
		return imageURL, nil
	}

	opts.GoogleAccessID = email
	opts.PrivateKey = []byte(key)

	signedURL, err := storage.SignedURL(bucketName, objectPath, opts)
	if err != nil {
		obs.Logger.ErrorContext(ctx, "error_generating_signed_url", "error", err, "bucket", bucketName, "object", objectPath)
		return imageURL, nil // Retornar URL original si falla
	}

	obs.Logger.InfoContext(ctx, "signed_url_generated", "original_url", imageURL, "signed_url", signedURL[:50]+"...")
	return signedURL, nil
}

// GenerateSignedWriteURL genera una signed URL de escritura (PUT) para subir una imagen a GCS
// Retorna la signed URL, la URL pública y el objectPath
func GenerateSignedWriteURL(ctx context.Context, obs observability.Observability, userID string, fileName string, contentType string) (uploadURL string, publicURL string, objectPath string, err error) {
	// Construir la ruta del objeto: userID/timestamp-random-filename
	timestamp := time.Now().Unix()
	randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
	objectPath = fmt.Sprintf("%s/%d-%s-%s", userID, timestamp, randomSuffix, fileName)

	bucketName := "micartapro-images"

	// Generar signed URL para PUT (subir)
	opts := &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: contentType,
	}

	// Obtener credenciales del service account
	credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsPath == "" {
		return "", "", "", fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	credsData, err := os.ReadFile(credsPath)
	if err != nil {
		return "", "", "", fmt.Errorf("error reading credentials: %w", err)
	}

	var credsJSON map[string]interface{}
	if err := json.Unmarshal(credsData, &credsJSON); err != nil {
		return "", "", "", fmt.Errorf("error parsing credentials: %w", err)
	}

	email, ok := credsJSON["client_email"].(string)
	if !ok || email == "" {
		return "", "", "", fmt.Errorf("missing client_email in credentials")
	}

	key, ok := credsJSON["private_key"].(string)
	if !ok || key == "" {
		return "", "", "", fmt.Errorf("missing private_key in credentials")
	}

	opts.GoogleAccessID = email
	opts.PrivateKey = []byte(key)

	uploadURL, err = storage.SignedURL(bucketName, objectPath, opts)
	if err != nil {
		return "", "", "", fmt.Errorf("error generating signed URL: %w", err)
	}

	// Construir la URL pública
	publicURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)

	obs.Logger.InfoContext(ctx, "signed_write_url_generated", "objectPath", objectPath, "userID", userID, "contentType", contentType)
	return uploadURL, publicURL, objectPath, nil
}
