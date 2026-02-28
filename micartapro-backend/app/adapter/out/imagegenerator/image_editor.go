package imagegenerator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
	supabase "github.com/supabase-community/supabase-go"
	"google.golang.org/genai"
)

type EditImage func(ctx context.Context, prompt string, referenceImageUrl string, aspectRatio string, imageCount int, menuItemId string, uploadURL string, publicURL string) (string, error)

func init() {
	ioc.Register(NewImageEditor)
}

func NewImageEditor(genaiClient *genai.Client, obs observability.Observability, supabaseClient *supabase.Client, gcsClient *storage.Client) (EditImage, error) {
	return func(ctx context.Context, prompt string, referenceImageUrl string, aspectRatio string, imageCount int, menuItemId string, uploadURL string, publicURL string) (string, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "edit_image")
		defer span.End()

		obs.Logger.InfoContext(spanCtx, "editing_image", "prompt", prompt, "referenceImageUrl", referenceImageUrl, "aspectRatio", aspectRatio, "imageCount", imageCount, "menuItemId", menuItemId)

		// 0. Validar que la imagen de referencia exista y tenga contenido antes de llamar a Vertex AI
		if err := ValidateReferenceImage(spanCtx, gcsClient, referenceImageUrl, obs); err != nil {
			return "", err
		}

		// 1. Generar signed URL si es necesario (para imágenes de GCS)
		// Esto evita descargar la imagen en el backend, ahorrando memoria RAM
		imageURL := referenceImageUrl
		if strings.Contains(referenceImageUrl, "storage.googleapis.com") {
			signedURL, err := GenerateSignedReadURL(spanCtx, gcsClient, obs, referenceImageUrl)
			if err != nil {
				obs.Logger.WarnContext(spanCtx, "error_generating_signed_url_fallback", "error", err, "url", referenceImageUrl)
				// Continuar con la URL original si falla
			} else {
				imageURL = signedURL
				obs.Logger.InfoContext(spanCtx, "using_signed_url_for_reference", "original_url", referenceImageUrl, "signed_url", signedURL[:50]+"...")
			}
		}

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

		// 4. Image-to-image con gemini-2.5-flash-image: GenerateContent con prompt + imagen de referencia
		// Parts: texto (prompt) + imagen inline, igual que en el ejemplo oficial
		imageModel := "gemini-2.5-flash-image"

		// Determinar el tamaño de imagen según el uso:
		// - "1K" para productos del catálogo (optimizado para apps móviles, carga más rápida)
		// - "2K" para imágenes de portada (mayor calidad visual)
		imageSize := "1K" // Por defecto optimizado para móviles
		if menuItemId == "cover" {
			imageSize = "2K" // Portadas usan mayor resolución
		}

		// Prefijo obligatorio: precisión sin inventar ni omitir ingredientes
		// Para sushi/piezas: respetar contenido de cada pieza Y envoltorio (Env) según la descripción
		precisionPrefix := "CRITICAL: Apply the requested changes while keeping EXACTLY the ingredients and elements visible. For sushi, pieces, or rolls: preserve each piece's exact filling (contenido) AND wrapper (envoltorio/Env). Do not add, invent, or omit anything. "
		optimizedPrompt := precisionPrefix + prompt
		if menuItemId != "cover" && menuItemId != "footer" {
			optimizedPrompt += " Optimize for mobile app product catalog: clear product visibility, professional food photography."
		}

		// 3. Usar FileData con URI en lugar de InlineData con bytes
		// Esto evita descargar la imagen en el backend, ahorrando memoria RAM
		// Gemini puede acceder directamente a la URL desde GCS
		parts := []*genai.Part{
			{Text: optimizedPrompt},
			{FileData: &genai.FileData{
				MIMEType: mimeType,
				FileURI:  imageURL,
			}},
		}
		contents := []*genai.Content{
			{
				Role:  "user",
				Parts: parts,
			},
		}
		if aspectRatio == "" {
			aspectRatio = "1:1"
		}
		config := &genai.GenerateContentConfig{
			ResponseModalities: []string{string(genai.ModalityText), string(genai.ModalityImage)},
			ImageConfig:        &genai.ImageConfig{AspectRatio: aspectRatio, ImageSize: imageSize},
		}

		obs.Logger.InfoContext(spanCtx, "calling_image_to_image", "model", imageModel, "imageUrl", imageURL, "mimeType", mimeType, "imageSize", imageSize, "menuItemId", menuItemId, "optimizedForMobile", menuItemId != "cover" && menuItemId != "footer", "usingFileData", true)
		respGen, err := genaiClient.Models.GenerateContent(spanCtx, imageModel, contents, config)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_image_to_image", "error", err, "prompt", prompt)
			return "", fmt.Errorf("error image-to-image: %w", err)
		}

		if len(respGen.Candidates) == 0 || respGen.Candidates[0].Content == nil {
			obs.Logger.ErrorContext(spanCtx, "no_candidates_or_content", "prompt", prompt)
			return "", fmt.Errorf("no image generated (empty response)")
		}

		// 5. Extraer bytes de la imagen de la respuesta (Parts con InlineData)
		var imgBytes []byte
		for _, p := range respGen.Candidates[0].Content.Parts {
			if p != nil && p.InlineData != nil && len(p.InlineData.Data) > 0 {
				imgBytes = p.InlineData.Data
				break
			}
		}
		if len(imgBytes) == 0 {
			obs.Logger.ErrorContext(spanCtx, "no_image_in_response", "prompt", prompt)
			return "", fmt.Errorf("no image bytes in model response")
		}

		obs.Logger.InfoContext(spanCtx, "image_to_image_success", "size_bytes", len(imgBytes))

		// 6. Subir usando la signed URL pre-firmada (con retry si ExpiredToken)
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

		// 7. Guardar en catalog_images
		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			return "", fmt.Errorf("userID is required but not found in context")
		}
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

		obs.Logger.InfoContext(spanCtx, "image_edited_and_uploaded", "publicURL", publicURL, "size_bytes", len(imgBytes))

		return publicURL, nil
	}, nil
}
