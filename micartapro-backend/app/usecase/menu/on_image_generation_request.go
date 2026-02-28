package menu

import (
	"context"
	"fmt"
	"micartapro/app/adapter/out/imagegenerator"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
)

type OnImageGenerationRequest func(ctx context.Context, input events.ImageGenerationRequestEvent) error

func init() {
	ioc.Register(NewOnImageGenerationRequest)
}

func NewOnImageGenerationRequest(
	obs observability.Observability,
	generateImage imagegenerator.GenerateImage,
	getUserCredits supabaserepo.GetUserCredits,
	consumeCredits supabaserepo.ConsumeCredits,
) OnImageGenerationRequest {
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

		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			return fmt.Errorf("invalid userID: %w", err)
		}
		userCredits, err := getUserCredits(spanCtx, parsedUserID)
		if err != nil {
			return fmt.Errorf("error getting user credits: %w", err)
		}
		if userCredits.Balance < billing.CreditsPerImageGeneration {
			return fmt.Errorf("insufficient credits: balance %d, required %d", userCredits.Balance, billing.CreditsPerImageGeneration)
		}
		sourceID := input.MenuID + ":" + input.MenuItemID + ":" + input.ImageType + ":gen"
		desc := "Generación de imagen"
		_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
			UserID:      parsedUserID,
			Amount:      billing.CreditsPerImageGeneration,
			Source:      "image.generation",
			SourceID:    &sourceID,
			Description: &desc,
		})
		if err != nil {
			if err == supabaserepo.ErrInsufficientCredits {
				return fmt.Errorf("insufficient credits for image generation")
			}
			return fmt.Errorf("error consuming credits: %w", err)
		}

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
