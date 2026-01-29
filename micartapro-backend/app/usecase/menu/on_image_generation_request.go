package menu

import (
	"context"
	"fmt"
	"micartapro/app/adapter/out/imagegenerator"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnImageGenerationRequest func(ctx context.Context, input events.ImageGenerationRequestEvent) error

func init() {
	ioc.Registry(
		NewOnImageGenerationRequest,
		observability.NewObservability,
		imagegenerator.NewImageGenerator,
	)
}

func NewOnImageGenerationRequest(
	obs observability.Observability,
	generateImage imagegenerator.GenerateImage) OnImageGenerationRequest {
	return func(ctx context.Context, input events.ImageGenerationRequestEvent) error {
		spanCtx, span := obs.Tracer.Start(ctx, "on_image_generation_request")
		defer span.End()

		// Verificar que userID esté en el contexto
		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			obs.Logger.ErrorContext(spanCtx, "userID_not_found_in_context",
				"menuId", input.MenuID,
				"menuItemId", input.MenuItemID,
			)
			return fmt.Errorf("userID is required but not found in context")
		}

		obs.Logger.InfoContext(spanCtx, "processing_image_generation_request",
			"menuId", input.MenuID,
			"menuItemId", input.MenuItemID,
			"imageType", input.ImageType,
			"publicURL", input.PublicURL,
			"userID", userID,
		)

		// Generar la imagen usando la signed URL pre-firmada
		// El generador ya guarda en catalog_images automáticamente
		publicURL, err := generateImage(spanCtx, input.Prompt, input.AspectRatio, input.ImageCount, input.UploadURL, input.PublicURL)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_generating_image", "error", err,
				"menuId", input.MenuID,
				"menuItemId", input.MenuItemID,
				"imageType", input.ImageType,
			)
			return fmt.Errorf("error generando imagen: %w", err)
		}

		obs.Logger.InfoContext(spanCtx, "image_generated_successfully",
			"menuId", input.MenuID,
			"menuItemId", input.MenuItemID,
			"imageType", input.ImageType,
			"publicURL", publicURL,
		)

		// La imagen ya fue guardada en catalog_images por el generador
		// La URL pública ya está asignada en el menú (placeholder)
		// No necesitamos actualizar el menú porque ya tiene la URL correcta

		return nil
	}
}
