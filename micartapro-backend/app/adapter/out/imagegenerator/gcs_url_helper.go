package imagegenerator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"micartapro/app/shared/infrastructure/observability"
)

// normalizeGCSURL corrige URLs de GCS mal formateadas (ej: "https.storage" → "https://storage")
// Es idempotente: puede aplicarse múltiples veces sin causar efectos secundarios
func normalizeGCSURL(url string) string {
	if url == "" {
		return url
	}
	
	// Verificar si la URL ya está correctamente formateada
	if strings.HasPrefix(url, "https://storage.googleapis.com") || strings.HasPrefix(url, "http://storage.googleapis.com") {
		return url
	}
	
	// Corregir "https.storage.googleapis.com" → "https://storage.googleapis.com"
	// Solo aplicar si el patrón incorrecto está presente
	if strings.Contains(url, "https.storage.googleapis.com") {
		url = strings.ReplaceAll(url, "https.storage.googleapis.com", "https://storage.googleapis.com")
	}
	if strings.Contains(url, "http.storage.googleapis.com") {
		url = strings.ReplaceAll(url, "http.storage.googleapis.com", "http://storage.googleapis.com")
	}
	
	// Manejar casos donde se duplicó "https" (httpshttps://storage...)
	url = strings.ReplaceAll(url, "httpshttps://", "https://")
	url = strings.ReplaceAll(url, "httphttp://", "http://")
	
	return url
}

// GenerateSignedReadURL genera una signed URL de lectura para una imagen en GCS usando el cliente de Storage.
// El cliente usa Application Default Credentials (en Cloud Run se detectan automáticamente).
// Si la URL no es de GCS, la retorna sin modificar.
func GenerateSignedReadURL(ctx context.Context, client *storage.Client, obs observability.Observability, imageURL string) (string, error) {
	// Normalizar URL primero (corregir "https.storage" → "https://storage")
	imageURL = normalizeGCSURL(imageURL)
	
	if client == nil {
		obs.Logger.WarnContext(ctx, "gcs_client_nil", "message", "using original URL")
		return imageURL, nil
	}
	if !strings.Contains(imageURL, "storage.googleapis.com") {
		obs.Logger.InfoContext(ctx, "url_not_gcs", "url", imageURL)
		return imageURL, nil
	}

	parts := strings.Split(imageURL, "storage.googleapis.com/")
	if len(parts) != 2 {
		obs.Logger.WarnContext(ctx, "invalid_gcs_url_format", "url", imageURL)
		return imageURL, nil
	}
	pathParts := strings.SplitN(parts[1], "/", 2)
	if len(pathParts) != 2 {
		obs.Logger.WarnContext(ctx, "invalid_gcs_path_format", "url", imageURL)
		return imageURL, nil
	}
	bucketName := pathParts[0]
	objectPath := pathParts[1]

	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(1 * time.Hour),
	}
	signedURL, err := client.Bucket(bucketName).SignedURL(objectPath, opts)
	if err != nil {
		obs.Logger.ErrorContext(ctx, "error_generating_signed_url", "error", err, "bucket", bucketName, "object", objectPath)
		return imageURL, nil
	}
	obs.Logger.InfoContext(ctx, "signed_url_generated", "original_url", imageURL, "signed_url", signedURL[:min(50, len(signedURL))]+"...")
	return signedURL, nil
}

// GenerateSignedWriteURL genera una signed URL de escritura (PUT) para subir una imagen a GCS usando el cliente de Storage.
// Retorna la signed URL, la URL pública y el objectPath.
func GenerateSignedWriteURL(ctx context.Context, client *storage.Client, obs observability.Observability, userID string, fileName string, contentType string) (uploadURL string, publicURL string, objectPath string, err error) {
	if client == nil {
		return "", "", "", fmt.Errorf("GCS client is nil")
	}
	timestamp := time.Now().Unix()
	randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
	objectPath = fmt.Sprintf("%s/%d-%s-%s", userID, timestamp, randomSuffix, fileName)
	bucketName := "micartapro-images"

	opts := &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: contentType,
	}
	uploadURL, err = client.Bucket(bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", "", "", fmt.Errorf("generating signed URL: %w", err)
	}
	publicURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
	obs.Logger.InfoContext(ctx, "signed_write_url_generated", "objectPath", objectPath, "userID", userID, "contentType", contentType)
	return uploadURL, publicURL, objectPath, nil
}
