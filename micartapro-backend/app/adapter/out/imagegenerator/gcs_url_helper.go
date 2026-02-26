package imagegenerator

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"micartapro/app/shared/infrastructure/observability"

	"cloud.google.com/go/storage"
)

// ErrReferenceImageNotAvailable se retorna cuando la imagen de referencia no existe o está vacía.
// Indica que el ítem no es candidato para edición image-to-image.
var ErrReferenceImageNotAvailable = fmt.Errorf("la imagen de referencia no está disponible o está vacía; no es candidata para edición")

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

// GenerateSignedWriteURLForReelScene genera signed URL para subir imagen de escena de reel.
// Path: {userID}/reel-scenes/{timestamp}-{random}-scene-{sceneIndex}.png
func GenerateSignedWriteURLForReelScene(ctx context.Context, client *storage.Client, obs observability.Observability, userID string, sceneIndex int, contentType string) (uploadURL string, publicURL string, objectPath string, err error) {
	return generateSignedWriteURLForReelMedia(ctx, client, obs, userID, sceneIndex, contentType, "png")
}

// GenerateSignedWriteURLForReelSceneVideo genera signed URL para subir video de escena de reel.
// Path: {userID}/reel-scenes/{timestamp}-{random}-scene-{sceneIndex}.mp4
func GenerateSignedWriteURLForReelSceneVideo(ctx context.Context, client *storage.Client, obs observability.Observability, userID string, sceneIndex int, contentType string) (uploadURL string, publicURL string, objectPath string, err error) {
	return generateSignedWriteURLForReelMedia(ctx, client, obs, userID, sceneIndex, contentType, "mp4")
}

func generateSignedWriteURLForReelMedia(ctx context.Context, client *storage.Client, obs observability.Observability, userID string, sceneIndex int, contentType string, ext string) (uploadURL string, publicURL string, objectPath string, err error) {
	if client == nil {
		return "", "", "", fmt.Errorf("GCS client is nil")
	}
	timestamp := time.Now().Unix()
	randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
	fileName := fmt.Sprintf("scene-%d.%s", sceneIndex, ext)
	objectPath = fmt.Sprintf("%s/reel-scenes/%d-%s-%s", userID, timestamp, randomSuffix, fileName)
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
	obs.Logger.InfoContext(ctx, "signed_write_url_reel_scene_generated", "objectPath", objectPath, "userID", userID, "sceneIndex", sceneIndex, "contentType", contentType)
	return uploadURL, publicURL, objectPath, nil
}

// RegenerateSignedWriteURL regenera una signed URL de escritura para el mismo objeto (publicURL).
// Útil cuando la URL original expiró (ExpiredToken) y se necesita reintentar el upload.
func RegenerateSignedWriteURL(ctx context.Context, client *storage.Client, obs observability.Observability, publicURL string, contentType string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("GCS client is nil")
	}
	url := normalizeGCSURL(publicURL)
	if !strings.Contains(url, "storage.googleapis.com/") {
		return "", fmt.Errorf("publicURL no es una URL de GCS: %s", publicURL)
	}
	parts := strings.Split(url, "storage.googleapis.com/")
	if len(parts) != 2 {
		return "", fmt.Errorf("formato de URL GCS inválido: %s", publicURL)
	}
	pathParts := strings.SplitN(parts[1], "/", 2)
	if len(pathParts) != 2 {
		return "", fmt.Errorf("path GCS inválido: %s", publicURL)
	}
	bucketName := pathParts[0]
	objectPath := pathParts[1]

	opts := &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: contentType,
	}
	uploadURL, err := client.Bucket(bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("regenerando signed URL: %w", err)
	}
	obs.Logger.InfoContext(ctx, "signed_write_url_regenerated", "objectPath", objectPath, "contentType", contentType)
	return uploadURL, nil
}

// ValidateReferenceImage verifica que la URL de referencia exista y tenga contenido antes de usarla en image-to-image.
// Evita errores URL_ERROR-ERROR_NOT_FOUND de Vertex AI cuando la URL apunta a un placeholder vacío.
func ValidateReferenceImage(ctx context.Context, client *storage.Client, imageURL string, obs observability.Observability) error {
	if imageURL == "" {
		return ErrReferenceImageNotAvailable
	}
	imageURL = normalizeGCSURL(imageURL)

	if strings.Contains(imageURL, "storage.googleapis.com/") && client != nil {
		parts := strings.Split(imageURL, "storage.googleapis.com/")
		if len(parts) != 2 {
			obs.Logger.WarnContext(ctx, "invalid_gcs_url_for_validation", "url", imageURL)
			return ErrReferenceImageNotAvailable
		}
		pathParts := strings.SplitN(parts[1], "/", 2)
		if len(pathParts) != 2 {
			return ErrReferenceImageNotAvailable
		}
		bucketName := pathParts[0]
		objectPath := pathParts[1]

		attrs, err := client.Bucket(bucketName).Object(objectPath).Attrs(ctx)
		if err != nil {
			if err == storage.ErrObjectNotExist {
				obs.Logger.InfoContext(ctx, "reference_image_not_found_gcs", "url", imageURL)
				return ErrReferenceImageNotAvailable
			}
			obs.Logger.WarnContext(ctx, "error_checking_reference_image", "error", err, "url", imageURL)
			return ErrReferenceImageNotAvailable
		}
		if attrs.Size == 0 {
			obs.Logger.InfoContext(ctx, "reference_image_empty", "url", imageURL)
			return ErrReferenceImageNotAvailable
		}
		return nil
	}

	// Para URLs no-GCS (signed, etc.): HEAD request
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, imageURL, nil)
	if err != nil {
		obs.Logger.WarnContext(ctx, "error_creating_head_request", "error", err, "url", imageURL)
		return ErrReferenceImageNotAvailable
	}
	clientHTTP := &http.Client{Timeout: 10 * time.Second}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		obs.Logger.WarnContext(ctx, "error_head_reference_image", "error", err, "url", imageURL)
		return ErrReferenceImageNotAvailable
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		obs.Logger.InfoContext(ctx, "reference_image_not_accessible", "url", imageURL, "status", resp.StatusCode)
		return ErrReferenceImageNotAvailable
	}
	if resp.ContentLength >= 0 && resp.ContentLength == 0 {
		obs.Logger.InfoContext(ctx, "reference_image_empty_http", "url", imageURL)
		return ErrReferenceImageNotAvailable
	}
	return nil
}
