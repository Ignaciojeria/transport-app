package imagegenerator

import (
	"context"
	"fmt"

	"micartapro/app/shared/infrastructure/ai"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

type GenerateImage func(ctx context.Context, prompt string, aspectRatio string, imageCount int) ([]byte, error)

func init() {
	ioc.Registry(NewImageGenerator, ai.NewClient, observability.NewObservability)
}

func NewImageGenerator(client *genai.Client, obs observability.Observability) (GenerateImage, error) {

	return func(ctx context.Context, prompt string, aspectRatio string, imageCount int) ([]byte, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "generate_image")
		defer span.End()

		obs.Logger.InfoContext(spanCtx, "generating_image", "prompt", prompt, "aspectRatio", aspectRatio, "imageCount", imageCount)

		// Usar el método GenerateImages del Models según models.go línea 5265
		config := &genai.GenerateImagesConfig{
			AspectRatio:    aspectRatio,
			NumberOfImages: 1,
		}

		resp, err := client.Models.GenerateImages(spanCtx, "imagen-4.0-ultra-generate-001", prompt, config)

		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_generating_image", "error", err, "prompt", prompt)
			return nil, fmt.Errorf("error generating image: %w", err)
		}

		if len(resp.GeneratedImages) == 0 {
			obs.Logger.ErrorContext(spanCtx, "no_images_generated", "prompt", prompt)
			return nil, fmt.Errorf("no images generated")
		}

		// Acceder a los bytes de la imagen: GeneratedImage.Image.ImageBytes
		generatedImage := resp.GeneratedImages[0]
		if generatedImage.Image == nil {
			obs.Logger.ErrorContext(spanCtx, "image_is_nil", "prompt", prompt)
			return nil, fmt.Errorf("generated image is nil")
		}

		imgBytes := generatedImage.Image.ImageBytes
		if len(imgBytes) == 0 {
			obs.Logger.ErrorContext(spanCtx, "image_bytes_empty", "prompt", prompt)
			return nil, fmt.Errorf("image bytes are empty")
		}

		obs.Logger.InfoContext(spanCtx, "image_generated_successfully", "size_bytes", len(imgBytes))

		return imgBytes, nil
	}, nil
}
