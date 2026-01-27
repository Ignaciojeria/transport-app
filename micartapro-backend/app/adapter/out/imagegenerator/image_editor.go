package imagegenerator

import (
	"context"
	"fmt"
	"io"
	"micartapro/app/shared/infrastructure/ai"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

type EditImage func(ctx context.Context, prompt string, referenceImageUrl string, aspectRatio string, imageCount int, menuItemId string) ([]byte, error)

func init() {
	ioc.Registry(NewImageEditor, ai.NewClient, observability.NewObservability)
}

func NewImageEditor(client *genai.Client, obs observability.Observability) (EditImage, error) {
	return func(ctx context.Context, prompt string, referenceImageUrl string, aspectRatio string, imageCount int, menuItemId string) ([]byte, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "edit_image")
		defer span.End()

		obs.Logger.InfoContext(spanCtx, "editing_image", "prompt", prompt, "referenceImageUrl", referenceImageUrl, "aspectRatio", aspectRatio, "imageCount", imageCount, "menuItemId", menuItemId)

		// 1. Descargar la imagen de referencia desde la URL
		obs.Logger.InfoContext(spanCtx, "downloading_reference_image", "url", referenceImageUrl)
		httpClient := &http.Client{}
		resp, err := httpClient.Get(referenceImageUrl)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_downloading_reference_image", "error", err, "url", referenceImageUrl)
			return nil, fmt.Errorf("error downloading reference image: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			obs.Logger.ErrorContext(spanCtx, "error_downloading_reference_image", "status", resp.StatusCode, "url", referenceImageUrl)
			return nil, fmt.Errorf("error downloading reference image: status %d", resp.StatusCode)
		}

		// 2. Leer los bytes de la imagen
		imageBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_reading_reference_image", "error", err)
			return nil, fmt.Errorf("error reading reference image: %w", err)
		}

		if len(imageBytes) == 0 {
			obs.Logger.ErrorContext(spanCtx, "reference_image_empty", "url", referenceImageUrl)
			return nil, fmt.Errorf("reference image is empty")
		}

		obs.Logger.InfoContext(spanCtx, "reference_image_downloaded", "size_bytes", len(imageBytes))

		// 3. Determinar MIME type basado en los primeros bytes de la imagen
		mimeType := "image/png"
		if len(imageBytes) >= 2 {
			if imageBytes[0] == 0xFF && imageBytes[1] == 0xD8 {
				mimeType = "image/jpeg"
			} else if len(imageBytes) >= 4 && imageBytes[0] == 0x89 && imageBytes[1] == 0x50 && imageBytes[2] == 0x4E && imageBytes[3] == 0x47 {
				mimeType = "image/png"
			} else if len(imageBytes) >= 4 && imageBytes[0] == 0x47 && imageBytes[1] == 0x49 && imageBytes[2] == 0x46 {
				mimeType = "image/gif"
			} else if len(imageBytes) >= 2 && imageBytes[0] == 0x42 && imageBytes[1] == 0x4D {
				mimeType = "image/bmp"
			}
		}

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
		
		// Agregar instrucciones de optimización móvil al prompt si es para productos
		optimizedPrompt := prompt
		if menuItemId != "cover" && menuItemId != "footer" {
			optimizedPrompt = prompt + " Optimize this image for mobile app product catalog display: ensure fast loading, clear product visibility, and professional food photography quality suitable for small screens."
		}
		
		parts := []*genai.Part{
			{Text: optimizedPrompt},
			{InlineData: &genai.Blob{Data: imageBytes, MIMEType: mimeType}},
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

		obs.Logger.InfoContext(spanCtx, "calling_image_to_image", "model", imageModel, "imageBytesSize", len(imageBytes), "mimeType", mimeType, "imageSize", imageSize, "menuItemId", menuItemId, "optimizedForMobile", menuItemId != "cover" && menuItemId != "footer")
		respGen, err := client.Models.GenerateContent(spanCtx, imageModel, contents, config)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_image_to_image", "error", err, "prompt", prompt)
			return nil, fmt.Errorf("error image-to-image: %w", err)
		}

		if len(respGen.Candidates) == 0 || respGen.Candidates[0].Content == nil {
			obs.Logger.ErrorContext(spanCtx, "no_candidates_or_content", "prompt", prompt)
			return nil, fmt.Errorf("no image generated (empty response)")
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
			return nil, fmt.Errorf("no image bytes in model response")
		}

		obs.Logger.InfoContext(spanCtx, "image_to_image_success", "size_bytes", len(imgBytes))

		return imgBytes, nil
	}, nil
}
