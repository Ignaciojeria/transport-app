package menu

import (
	"context"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnMenuCreateRequest func(ctx context.Context, input events.MenuCreateRequest) error

func init() {
	ioc.Registry(NewOnMenuCreateRequest,
		observability.NewObservability,
		storage.NewSaveMenu,
		supabaserepo.NewSaveMenu,
		eventprocessing.NewPublisherStrategy,
	)
}

func NewOnMenuCreateRequest(
	obs observability.Observability,
	saveMenuStorage storage.SaveMenu,
	saveMenuSupabase supabaserepo.SaveMenu,
	publisherManager eventprocessing.PublisherManager) OnMenuCreateRequest {
	return func(ctx context.Context, input events.MenuCreateRequest) error {
		obs.Logger.InfoContext(ctx, "on_menu_create_request", "menuId", input.ID)
		spanCtx, span := obs.Tracer.Start(ctx, "on_menu_create_request")
		defer span.End()

		// Guardar referencias a las solicitudes de generación/edición antes de limpiar
		// para poder publicar eventos después de guardar el menú
		coverGenReq := input.CoverImageGenerationRequest
		coverEditReq := input.CoverImageEditionRequest
		imageGenReqs := make([]events.ImageGenerationRequest, len(input.ImageGenerationRequests))
		copy(imageGenReqs, input.ImageGenerationRequests)
		imageEditReqs := make([]events.ImageEditionRequest, len(input.ImageEditionRequests))
		copy(imageEditReqs, input.ImageEditionRequests)

		// Limpiar campos temporales de solicitudes de generación/edición de imágenes
		// antes de guardar (estos campos no deben persistirse)
		input.Clean()

		// Normalizar URLs de imágenes antes de guardar (corrige httpshttps, https.storage, etc.)
		input.NormalizeImageURLs()

		// Guardar en GCS (storage)
		err := saveMenuStorage(spanCtx, input)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_saving_menu_to_storage", "error", err)
			return err
		}

		// Guardar en Supabase
		err = saveMenuSupabase(spanCtx, input)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_saving_menu_to_supabase", "error", err)
			return err
		}

		obs.Logger.InfoContext(spanCtx, "menu_saved_successfully", "menuId", input.ID)

		// Publicar eventos de generación/edición de imágenes de forma asíncrona
		// después de que el menú haya sido guardado exitosamente
		if err := publishImageGenerationEvents(spanCtx, obs, publisherManager, input.ID, coverGenReq, coverEditReq, imageGenReqs, imageEditReqs); err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_publishing_image_generation_events", "error", err)
			// No retornar error aquí porque el menú ya fue guardado exitosamente
			// Los eventos de imágenes se pueden reintentar más tarde
		}

		return nil
	}
}

// publishImageGenerationEvents publica eventos individuales para generar/editar imágenes de forma asíncrona
func publishImageGenerationEvents(
	ctx context.Context,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager,
	menuID string,
	coverGenReq *events.CoverImageGenerationRequest,
	coverEditReq *events.CoverImageEditionRequest,
	imageGenReqs []events.ImageGenerationRequest,
	imageEditReqs []events.ImageEditionRequest,
) error {
	// Publicar evento para generar imagen de portada
	if coverGenReq != nil && coverGenReq.UploadURL != "" && coverGenReq.PublicURL != "" {
		event := events.ImageGenerationRequestEvent{
			MenuID:      menuID,
			Prompt:      coverGenReq.Prompt,
			AspectRatio: "16:9",
			ImageCount:  coverGenReq.ImageCount,
			UploadURL:   coverGenReq.UploadURL,
			PublicURL:   coverGenReq.PublicURL,
			ImageType:   "cover",
		}
		if err := publisherManager.Publish(ctx, eventprocessing.PublishRequest{
			Topic:       "micartapro.events",
			Source:      "micartapro.menu.create",
			OrderingKey: menuID,
			Event:       event,
		}); err != nil {
			obs.Logger.ErrorContext(ctx, "error_publishing_cover_generation_event", "error", err)
			return err
		}
		obs.Logger.InfoContext(ctx, "cover_generation_event_published", "menuId", menuID, "publicURL", coverGenReq.PublicURL)
	}

	// Publicar evento para editar imagen de portada
	if coverEditReq != nil && coverEditReq.UploadURL != "" && coverEditReq.PublicURL != "" {
		event := events.ImageEditionRequestEvent{
			MenuID:            menuID,
			Prompt:           coverEditReq.Prompt,
			ReferenceImageUrl: coverEditReq.ReferenceImageUrl,
			AspectRatio:      "16:9",
			ImageCount:       coverEditReq.ImageCount,
			UploadURL:        coverEditReq.UploadURL,
			PublicURL:        coverEditReq.PublicURL,
			ImageType:        "cover",
		}
		if err := publisherManager.Publish(ctx, eventprocessing.PublishRequest{
			Topic:       "micartapro.events",
			Source:      "micartapro.menu.create",
			OrderingKey: menuID,
			Event:       event,
		}); err != nil {
			obs.Logger.ErrorContext(ctx, "error_publishing_cover_edition_event", "error", err)
			return err
		}
		obs.Logger.InfoContext(ctx, "cover_edition_event_published", "menuId", menuID, "publicURL", coverEditReq.PublicURL)
	}

	// Publicar eventos para generar imágenes de items (deduplicar por menuItemId para evitar duplicados si el agente repite)
	seenGen := make(map[string]bool)
	for _, req := range imageGenReqs {
		if req.UploadURL == "" || req.PublicURL == "" {
			continue
		}
		if seenGen[req.MenuItemID] {
			obs.Logger.WarnContext(ctx, "skipping_duplicate_image_generation_request", "menuItemId", req.MenuItemID)
			continue
		}
		seenGen[req.MenuItemID] = true
		event := events.ImageGenerationRequestEvent{
			MenuID:      menuID,
			MenuItemID:  req.MenuItemID,
			Prompt:      req.Prompt,
			AspectRatio: req.AspectRatio,
			ImageCount:  req.ImageCount,
			UploadURL:   req.UploadURL,
			PublicURL:   req.PublicURL,
			ImageType:   "item",
		}
		if err := publisherManager.Publish(ctx, eventprocessing.PublishRequest{
			Topic:       "micartapro.events",
			Source:      "micartapro.menu.create",
			OrderingKey: menuID,
			Event:       event,
		}); err != nil {
			obs.Logger.ErrorContext(ctx, "error_publishing_item_generation_event", "error", err, "menuItemId", req.MenuItemID)
			continue // Continuar con otros items aunque uno falle
		}
		obs.Logger.InfoContext(ctx, "item_generation_event_published", "menuId", menuID, "menuItemId", req.MenuItemID, "publicURL", req.PublicURL)
	}

	// Publicar eventos para editar imágenes de items (deduplicar por menuItemId para evitar duplicados)
	seenEdit := make(map[string]bool)
	for _, req := range imageEditReqs {
		if req.UploadURL == "" || req.PublicURL == "" {
			continue
		}
		if seenEdit[req.MenuItemID] {
			obs.Logger.WarnContext(ctx, "skipping_duplicate_image_edition_request", "menuItemId", req.MenuItemID)
			continue
		}
		seenEdit[req.MenuItemID] = true
		event := events.ImageEditionRequestEvent{
			MenuID:            menuID,
			MenuItemID:        req.MenuItemID,
			Prompt:           req.Prompt,
			ReferenceImageUrl: req.ReferenceImageUrl,
			AspectRatio:      req.AspectRatio,
			ImageCount:       req.ImageCount,
			UploadURL:        req.UploadURL,
			PublicURL:        req.PublicURL,
			ImageType:        "item",
		}
		if err := publisherManager.Publish(ctx, eventprocessing.PublishRequest{
			Topic:       "micartapro.events",
			Source:      "micartapro.menu.create",
			OrderingKey: menuID,
			Event:       event,
		}); err != nil {
			obs.Logger.ErrorContext(ctx, "error_publishing_item_edition_event", "error", err, "menuItemId", req.MenuItemID)
			continue // Continuar con otros items aunque uno falle
		}
		obs.Logger.InfoContext(ctx, "item_edition_event_published", "menuId", menuID, "menuItemId", req.MenuItemID, "publicURL", req.PublicURL)
	}

	return nil
}
