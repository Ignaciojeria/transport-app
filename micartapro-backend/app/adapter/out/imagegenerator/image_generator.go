package imagegenerator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	"cloud.google.com/go/storage"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
	supabase "github.com/supabase-community/supabase-go"
	"google.golang.org/genai"
)

type GenerateImage func(ctx context.Context, prompt string, aspectRatio string, imageCount int, uploadURL string, publicURL string) (string, error)

func init() {
	ioc.Register(NewImageGenerator)
}

func NewImageGenerator(client *genai.Client, obs observability.Observability, supabaseClient *supabase.Client, gcsClient *storage.Client) (GenerateImage, error) {
	return func(ctx context.Context, prompt string, aspectRatio string, imageCount int, uploadURL string, publicURL string) (string, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "generate_image")
		defer span.End()

		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			return "", fmt.Errorf("userID is required but not found in context")
		}

		obs.Logger.InfoContext(spanCtx, "generating_image", "prompt", prompt, "aspectRatio", aspectRatio, "imageCount", imageCount)

		// Prefijo obligatorio: generar exactamente lo descrito, sin inventar ni omitir
		// Para sushi/piezas: respetar contenido de cada pieza Y envoltorio (Env) según la descripción
		precisePrompt := "CRITICAL: Generate exactly and only what is described below. Do not add, invent, or omit any ingredients, elements, or details. For sushi, pieces, or rolls: each piece must show its exact filling (contenido) AND wrapper (envoltorio/Env) as stated. The image must show precisely what is stated—nothing more, nothing less. Professional food photography style.\n\n" + prompt

		// 1. Generar la imagen con Gemini
		config := &genai.GenerateImagesConfig{
			AspectRatio:    aspectRatio,
			NumberOfImages: 1,
		}

		resp, err := client.Models.GenerateImages(spanCtx, "imagen-4.0-ultra-generate-001", precisePrompt, config)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_generating_image", "error", err, "prompt", prompt)
			return "", fmt.Errorf("error generating image: %w", err)
		}

		if len(resp.GeneratedImages) == 0 {
			obs.Logger.ErrorContext(spanCtx, "no_images_generated", "prompt", prompt)
			return "", fmt.Errorf("no images generated")
		}

		generatedImage := resp.GeneratedImages[0]
		if generatedImage.Image == nil {
			obs.Logger.ErrorContext(spanCtx, "image_is_nil", "prompt", prompt)
			return "", fmt.Errorf("generated image is nil")
		}

		imgBytes := generatedImage.Image.ImageBytes
		if len(imgBytes) == 0 {
			obs.Logger.ErrorContext(spanCtx, "image_bytes_empty", "prompt", prompt)
			return "", fmt.Errorf("image bytes are empty")
		}

		obs.Logger.InfoContext(spanCtx, "image_generated_successfully", "size_bytes", len(imgBytes))

		// 2. Determinar MIME type basado en la URL pública de destino (publicURL)
		// Esto es crítico porque debe coincidir con el ContentType usado al generar la signed URL
		mimeType := "image/png" // Default (la mayoría de las imágenes generadas son PNG)
		lowerPublicURL := strings.ToLower(publicURL)
		if strings.Contains(lowerPublicURL, ".png") {
			mimeType = "image/png"
		} else if strings.Contains(lowerPublicURL, ".jpg") || strings.Contains(lowerPublicURL, ".jpeg") {
			mimeType = "image/jpeg"
		} else if strings.Contains(lowerPublicURL, ".gif") {
			mimeType = "image/gif"
		} else if strings.Contains(lowerPublicURL, ".webp") {
			mimeType = "image/webp"
		}

		obs.Logger.InfoContext(spanCtx, "mime_type_determined", "mimeType", mimeType, "publicURL", publicURL)

		// 3. Subir usando la signed URL pre-firmada (con retry si ExpiredToken)
		uploadURLToUse := uploadURL
		for attempt := 0; attempt < 2; attempt++ {
			req, err := http.NewRequestWithContext(spanCtx, "PUT", uploadURLToUse, bytes.NewReader(imgBytes))
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_creating_upload_request", "error", err)
				return "", fmt.Errorf("error creating upload request: %w", err)
			}
			req.Header.Set("Content-Type", mimeType)

			httpClient := &http.Client{Timeout: 30 * time.Second}
			respUpload, err := httpClient.Do(req)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_uploading_image", "error", err)
				return "", fmt.Errorf("error uploading image: %w", err)
			}
			body, _ := io.ReadAll(respUpload.Body)
			respUpload.Body.Close()

			if respUpload.StatusCode >= 200 && respUpload.StatusCode < 300 {
				obs.Logger.InfoContext(spanCtx, "image_uploaded_successfully", "publicURL", publicURL, "size_bytes", len(imgBytes))
				break
			}

			// Si ExpiredToken y tenemos gcsClient, regenerar URL y reintentar
			if attempt == 0 && respUpload.StatusCode == 400 && strings.Contains(string(body), "ExpiredToken") && gcsClient != nil {
				newURL, regenErr := RegenerateSignedWriteURL(spanCtx, gcsClient, obs, publicURL, mimeType)
				if regenErr != nil {
					obs.Logger.WarnContext(spanCtx, "error_regenerating_signed_url", "error", regenErr)
				} else {
					obs.Logger.InfoContext(spanCtx, "retrying_upload_with_regenerated_url", "publicURL", publicURL)
					uploadURLToUse = newURL
					continue
				}
			}

			obs.Logger.ErrorContext(spanCtx, "error_upload_status", "status", respUpload.StatusCode, "response_body", string(body), "publicURL", publicURL, "mimeType", mimeType, "content_length", len(imgBytes))
			return "", fmt.Errorf("error uploading image: status %d: %s", respUpload.StatusCode, string(body))
		}

		// 4. Guardar en catalog_images
		record := map[string]interface{}{
			"id":        uuid.New().String(),
			"image_url": publicURL,
			"user_id":   userID,
		}

		_, _, err = supabaseClient.From("catalog_images").
			Insert(record, false, "", "", "").
			Execute()

		if err != nil {
			obs.Logger.WarnContext(spanCtx, "error_saving_to_catalog_images", "error", err, "publicURL", publicURL, "userID", userID)
		} else {
			obs.Logger.InfoContext(spanCtx, "image_saved_to_catalog_images", "publicURL", publicURL, "userID", userID)
		}

		return publicURL, nil
	}, nil
}
