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

type OnImageEditionRequest func(ctx context.Context, input events.ImageEditionRequestEvent) error

func init() {
	ioc.Register(NewOnImageEditionRequest)
}

func NewOnImageEditionRequest(
	obs observability.Observability,
	editImage imagegenerator.EditImage,
	getUserCredits supabaserepo.GetUserCredits,
	consumeCredits supabaserepo.ConsumeCredits,
) OnImageEditionRequest {
	return func(ctx context.Context, input events.ImageEditionRequestEvent) error {
		spanCtx, span := obs.Tracer.Start(ctx, "on_image_edition_request")
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

		obs.Logger.InfoContext(spanCtx, "processing_image_edition_request",
			"menuId", input.MenuID,
			"menuItemId", input.MenuItemID,
			"imageType", input.ImageType,
			"publicURL", input.PublicURL,
			"referenceImageUrl", input.ReferenceImageUrl,
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
		if userCredits.Balance < billing.CreditsPerImageEdition {
			return fmt.Errorf("insufficient credits: balance %d, required %d", userCredits.Balance, billing.CreditsPerImageEdition)
		}
		sourceID := input.MenuID + ":" + input.MenuItemID + ":" + input.ImageType + ":edit"
		desc := "Edición de imagen"
		_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
			UserID:      parsedUserID,
			Amount:      billing.CreditsPerImageEdition,
			Source:      "image.edition",
			SourceID:    &sourceID,
			Description: &desc,
		})
		if err != nil {
			if err == supabaserepo.ErrInsufficientCredits {
				return fmt.Errorf("insufficient credits for image edition")
			}
			return fmt.Errorf("error consuming credits: %w", err)
		}

		// Determinar menuItemId para el editor (vacío para cover)
		menuItemId := input.MenuItemID
		if input.ImageType == "cover" {
			menuItemId = "cover"
		}

		// Editar la imagen usando la signed URL pre-firmada
		// El editor ya guarda en catalog_images automáticamente
		publicURL, err := editImage(spanCtx, input.Prompt, input.ReferenceImageUrl, input.AspectRatio, input.ImageCount, menuItemId, input.UploadURL, input.PublicURL)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_editing_image", "error", err,
				"menuId", input.MenuID,
				"menuItemId", input.MenuItemID,
				"imageType", input.ImageType,
			)
			return fmt.Errorf("error editando imagen: %w", err)
		}

		obs.Logger.InfoContext(spanCtx, "image_edited_successfully",
			"menuId", input.MenuID,
			"menuItemId", input.MenuItemID,
			"imageType", input.ImageType,
			"publicURL", publicURL,
		)

		// La imagen ya fue guardada en catalog_images por el editor
		// La URL pública ya está asignada en el menú (placeholder)
		// No necesitamos actualizar el menú porque ya tiene la URL correcta

		return nil
	}
}
