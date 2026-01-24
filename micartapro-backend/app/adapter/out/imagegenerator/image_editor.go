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

		// 3. Crear Image con los bytes descargados
		// Determinar el MIME type basado en los primeros bytes de la imagen
		mimeType := "image/png" // Por defecto PNG
		if len(imageBytes) >= 2 {
			// Verificar magic numbers para diferentes formatos
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

		// Validar que tenemos bytes válidos antes de crear la imagen
		if len(imageBytes) == 0 {
			obs.Logger.ErrorContext(spanCtx, "image_bytes_empty_after_download", "url", referenceImageUrl)
			return nil, fmt.Errorf("image bytes are empty after download")
		}

		referenceImageData := &genai.Image{
			ImageBytes: imageBytes,
			MIMEType:   mimeType,
		}

		// Validar que la imagen tiene bytes antes de continuar
		if len(referenceImageData.ImageBytes) == 0 {
			obs.Logger.ErrorContext(spanCtx, "reference_image_data_empty", "url", referenceImageUrl)
			return nil, fmt.Errorf("reference image data is empty")
		}

		// 4. Crear ReferenceImage usando NewRawReferenceImage
		// El referenceID debe ser un valor válido (1 o mayor), no 0
		// Usamos 1 como referenceID para vincular explícitamente la imagen
		referenceID := int32(1)
		referenceImage := genai.NewRawReferenceImage(referenceImageData, referenceID)

		// Validar que la referencia se creó correctamente
		if referenceImage == nil {
			obs.Logger.ErrorContext(spanCtx, "reference_image_creation_failed", "url", referenceImageUrl)
			return nil, fmt.Errorf("failed to create reference image")
		}

		// 5. Configurar EditImageConfig
		// Asegurarse de que la configuración esté completa
		config := &genai.EditImageConfig{
			AspectRatio:    aspectRatio,
			NumberOfImages: 1,
		}

		// 6. Crear array de ReferenceImage
		referenceImages := []genai.ReferenceImage{referenceImage}

		// 7. Usar el prompt del usuario tal cual, sin modificaciones
		refinedPrompt := prompt

		// 8. Llamar a EditImage
		// IMPORTANTE: Imagen 4 no tiene modelo de edición, debemos usar imagen-3.0-capability-001
		// Usar un modelo de generación (imagen-4.0-ultra-generate-001) causa el error:
		// "No uri or raw bytes are provided in media content"
		editModel := "imagen-3.0-capability-001"
		obs.Logger.InfoContext(spanCtx, "calling_edit_image", "model", editModel, "referenceID", referenceID, "imageBytesSize", len(imageBytes), "mimeType", mimeType)
		respEdit, err := client.Models.EditImage(spanCtx, editModel, refinedPrompt, referenceImages, config)
		if err != nil {
			// Mejorar el manejo de errores para obtener más detalles
			obs.Logger.ErrorContext(spanCtx, "error_editing_image",
				"error", err,
				"prompt", prompt,
				"referenceID", referenceID,
				"imageBytesSize", len(imageBytes),
				"mimeType", mimeType,
				"referenceImagesCount", len(referenceImages))

			// Intentar extraer más detalles del error si es posible
			if err.Error() != "" {
				obs.Logger.ErrorContext(spanCtx, "detailed_error_message", "errorDetails", err.Error())
			}

			return nil, fmt.Errorf("error editing image: %w", err)
		}

		if len(respEdit.GeneratedImages) == 0 {
			obs.Logger.ErrorContext(spanCtx, "no_images_edited", "prompt", prompt)
			return nil, fmt.Errorf("no images edited")
		}

		// 6. Acceder a los bytes de la imagen editada
		editedImage := respEdit.GeneratedImages[0]
		if editedImage.Image == nil {
			obs.Logger.ErrorContext(spanCtx, "edited_image_is_nil", "prompt", prompt)
			return nil, fmt.Errorf("edited image is nil")
		}

		imgBytes := editedImage.Image.ImageBytes
		if len(imgBytes) == 0 {
			obs.Logger.ErrorContext(spanCtx, "edited_image_bytes_empty", "prompt", prompt)
			return nil, fmt.Errorf("edited image bytes are empty")
		}

		obs.Logger.InfoContext(spanCtx, "image_edited_successfully", "size_bytes", len(imgBytes))

		return imgBytes, nil
	}, nil
}
