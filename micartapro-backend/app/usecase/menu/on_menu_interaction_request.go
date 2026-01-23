package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"micartapro/app/adapter/out/agents"
	"micartapro/app/adapter/out/imagegenerator"
	"micartapro/app/adapter/out/imageuploader"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
)

type OnMenuInteractionRequest func(ctx context.Context, input events.MenuInteractionRequest) (string, error)

func init() {
	ioc.Registry(
		NewOnMenuInteractionRequest,
		observability.NewObservability,
		agents.NewMenuInteractionAgent,
		eventprocessing.NewPublisherStrategy,
		supabaserepo.NewGetMenuById,
		imagegenerator.NewImageGenerator,
		imageuploader.NewImageUploader)
}

func NewOnMenuInteractionRequest(
	obs observability.Observability,
	menuInteractionAgent agents.MenuInteractionAgent,
	publisherManager eventprocessing.PublisherManager,
	getMenuById supabaserepo.GetMenuById,
	generateImage imagegenerator.GenerateImage,
	uploadImage imageuploader.UploadImage) OnMenuInteractionRequest {
	return func(ctx context.Context, input events.MenuInteractionRequest) (string, error) {

		// 1. Verificar y propagar el versionID del contexto para guardar la nueva versión generada
		// El versionID se usa cuando se guarda el menú, NO para leer el menú actual
		if versionID, ok := sharedcontext.VersionIDFromContext(ctx); ok && versionID != "" {
			obs.Logger.InfoContext(ctx, "versionID from context will be used when saving new menu version", "menuID", input.MenuID, "versionID", versionID)
		} else {
			obs.Logger.InfoContext(ctx, "no versionID in context, a new one will be generated when saving", "menuID", input.MenuID)
		}

		// 2. Obtener el contenido del menú desde Supabase usando el menuID
		// IMPORTANTE: Pasamos string vacío para que getMenuById SIEMPRE use current_version_id
		// NO usamos el versionID del contexto para leer, solo para guardar después
		obs.Logger.InfoContext(ctx, "getting menu content using current_version_id from menu (ignoring versionID from context)", "menuID", input.MenuID)

		// 3. Obtener el contenido del menú desde menu_versions en Supabase
		// Esta función obtiene el current_version_id del menú desde la tabla menus
		// y luego busca el contenido en menu_versions usando ese current_version_id
		// El contenido obtenido será el menú actual que se pasará al prompt
		menu, err := getMenuById(ctx, input.MenuID, "")
		if err != nil && err != supabaserepo.ErrMenuNotFound {
			obs.Logger.ErrorContext(ctx, "error getting menu from supabase", "error", err, "menuID", input.MenuID)
			return "", err
		}
		
		// Si el menú no se encuentra, continuamos con un menú vacío (comportamiento original)
		if err == supabaserepo.ErrMenuNotFound {
			obs.Logger.WarnContext(ctx, "menu not found in supabase, continuing with empty menu", "menuID", input.MenuID)
		} else {
			obs.Logger.InfoContext(ctx, "menu content retrieved from menu_versions using current_version_id, will be passed to prompt", "menuID", input.MenuID)
		}
		
		// 4. Asignar el contenido del menú obtenido desde menu_versions (usando current_version_id) al request
		// Este contenido será convertido a toon format y pasado al prompt como MENU_ACTUAL
		// Representa la versión actual del menú que el agente debe conocer para aplicar los cambios
		// NOTA: El versionID del contexto se mantiene y se usará cuando se guarde el menú modificado
		input.JsonMenu = menu

		// 5. Llamar al Agente Procesador (que hace la inspección del FunctionCall)
		// El contexto con el versionID se propaga automáticamente y se usará al guardar el menú
		agentResp, err := menuInteractionAgent(ctx, input)
		if err != nil {
			return "", err
		}

		// 2. Lógica Limpia: ¿Es texto o un comando?

		if agentResp.TextResponse != "" {
			// Caso 1: Texto conversacional (Ej. "Faltan los precios")
			return agentResp.TextResponse, nil
		}

		if agentResp.CommandName == "createMenu" {
			// Caso 2: Comando. Hacemos el mapeo (la única deserialización necesaria)
			var createRequest events.MenuCreateRequest

			// Mapeo directo del JSON crudo a tu objeto de dominio
			if err := json.Unmarshal(agentResp.CommandArgs, &createRequest); err != nil {
				return "", fmt.Errorf("error al mapear a MenuCreateRequest: %w", err)
			}
			createRequest.ID = input.MenuID

			// Procesar imagen de portada si hay solicitud
			if createRequest.CoverImageGenerationRequest != nil {
				obs.Logger.InfoContext(ctx, "processing_cover_image_generation_request")
				
				if err := processCoverImageGenerationRequest(ctx, obs, &createRequest, generateImage, uploadImage); err != nil {
					obs.Logger.ErrorContext(ctx, "error_processing_cover_image_generation", "error", err)
					return "", fmt.Errorf("error al procesar generación de imagen de portada: %w", err)
				}
				
				obs.Logger.InfoContext(ctx, "cover_image_generation_completed")
			}

			// Procesar imágenes si hay solicitudes de generación
			if len(createRequest.ImageGenerationRequests) > 0 {
				obs.Logger.InfoContext(ctx, "processing_image_generation_requests", "count", len(createRequest.ImageGenerationRequests))
				
				if err := processImageGenerationRequests(ctx, obs, &createRequest, generateImage, uploadImage); err != nil {
					obs.Logger.ErrorContext(ctx, "error_processing_image_generation", "error", err)
					return "", fmt.Errorf("error al procesar generación de imágenes: %w", err)
				}
				
				obs.Logger.InfoContext(ctx, "image_generation_completed", "count", len(createRequest.ImageGenerationRequests))
			}

			if err := publisherManager.Publish(ctx, eventprocessing.PublishRequest{
				Topic:  "micartapro.events",
				Source: "micartapro.agent.menu.interaction",
				Event:  createRequest,
			}); err != nil {
				return "", fmt.Errorf("error al publicar evento de creación de menú: %w", err)
			}

			return "¡Menú creado exitosamente! Se ha notificado a todos los servicios.", nil
		}

		// ... manejar otros comandos

		return "Lo siento, no pude entender la acción solicitada.", nil
	}
}

// processCoverImageGenerationRequest procesa la solicitud de generación de imagen de portada
func processCoverImageGenerationRequest(
	ctx context.Context,
	obs observability.Observability,
	createRequest *events.MenuCreateRequest,
	generateImage imagegenerator.GenerateImage,
	uploadImage imageuploader.UploadImage,
) error {
	spanCtx, span := obs.Tracer.Start(ctx, "process_cover_image")
	defer span.End()

	coverReq := createRequest.CoverImageGenerationRequest
	obs.Logger.InfoContext(spanCtx, "generating_cover_image", "prompt", coverReq.Prompt)

	// 1. Generar la imagen con aspect ratio 16:9 (horizontalmente larga y verticalmente corta, tipo foto portada LinkedIn)
	// Nota: La API solo acepta aspect ratios predefinidos, 16:9 es el más ancho disponible
	imageBytes, err := generateImage(spanCtx, coverReq.Prompt, "16:9", coverReq.ImageCount)
	if err != nil {
		obs.Logger.ErrorContext(spanCtx, "error_generating_cover_image", "error", err)
		return fmt.Errorf("error generando imagen de portada: %w", err)
	}

	// 2. Crear nombre de archivo único
	fileName := fmt.Sprintf("cover-%s.png", uuid.New().String()[:8])

	// 3. Subir la imagen a Supabase Storage
	publicURL, err := uploadImage(spanCtx, imageBytes, fileName)
	if err != nil {
		obs.Logger.ErrorContext(spanCtx, "error_uploading_cover_image", "error", err)
		return fmt.Errorf("error subiendo imagen de portada: %w", err)
	}

	// 4. Asignar la URL directamente al coverImage
	createRequest.CoverImage = publicURL
	obs.Logger.InfoContext(spanCtx, "cover_image_processed_successfully", "publicURL", publicURL)

	return nil
}

// processImageGenerationRequests procesa todas las solicitudes de generación de imágenes en paralelo
func processImageGenerationRequests(
	ctx context.Context,
	obs observability.Observability,
	createRequest *events.MenuCreateRequest,
	generateImage imagegenerator.GenerateImage,
	uploadImage imageuploader.UploadImage,
) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(createRequest.ImageGenerationRequests))
	
	// Mapa para almacenar las URLs generadas por menuItemId
	imageURLs := make(map[string]string)
	var mu sync.Mutex

	for _, imgReq := range createRequest.ImageGenerationRequests {
		wg.Add(1)
		go func(req events.ImageGenerationRequest) {
			defer wg.Done()

			spanCtx, span := obs.Tracer.Start(ctx, "process_single_image")
			defer span.End()

			obs.Logger.InfoContext(spanCtx, "generating_image_for_item", "menuItemId", req.MenuItemID, "prompt", req.Prompt)

			// 1. Generar la imagen
			imageBytes, err := generateImage(spanCtx, req.Prompt, req.AspectRatio, req.ImageCount)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_generating_image", "error", err, "menuItemId", req.MenuItemID)
				errChan <- fmt.Errorf("error generando imagen para %s: %w", req.MenuItemID, err)
				return
			}

			// 2. Crear nombre de archivo único
			fileName := fmt.Sprintf("%s-%s.png", req.MenuItemID, uuid.New().String()[:8])

			// 3. Subir la imagen a Supabase Storage
			publicURL, err := uploadImage(spanCtx, imageBytes, fileName)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_uploading_image", "error", err, "menuItemId", req.MenuItemID)
				errChan <- fmt.Errorf("error subiendo imagen para %s: %w", req.MenuItemID, err)
				return
			}

			// 4. Guardar la URL en el mapa
			mu.Lock()
			imageURLs[req.MenuItemID] = publicURL
			mu.Unlock()

			obs.Logger.InfoContext(spanCtx, "image_processed_successfully", "menuItemId", req.MenuItemID, "publicURL", publicURL)
		}(imgReq)
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()
	close(errChan)

	// Verificar si hubo errores
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// 5. Actualizar el menú con las URLs de las imágenes
	updateMenuWithImageURLs(createRequest, imageURLs, obs, ctx)

	return nil
}

// updateMenuWithImageURLs actualiza el menú asignando las URLs de las imágenes a los items/sides correspondientes
// También maneja IDs especiales: "footer" para footerImage (cover ahora se maneja por separado)
func updateMenuWithImageURLs(
	createRequest *events.MenuCreateRequest,
	imageURLs map[string]string,
	obs observability.Observability,
	ctx context.Context,
) {
	for menuItemID, imageURL := range imageURLs {
		// Manejar IDs especiales para imágenes del menú
		if menuItemID == "footer" {
			createRequest.FooterImage = imageURL
			obs.Logger.InfoContext(ctx, "image_url_assigned_to_footer", "imageURL", imageURL)
			continue
		}
		
		// Buscar el item o side en el menú
		found := false
		
		// Buscar en items
		for i := range createRequest.Menu {
			for j := range createRequest.Menu[i].Items {
				if createRequest.Menu[i].Items[j].ID == menuItemID {
					createRequest.Menu[i].Items[j].PhotoUrl = imageURL
					obs.Logger.InfoContext(ctx, "image_url_assigned_to_item", "menuItemId", menuItemID, "imageURL", imageURL)
					found = true
					break
				}
				
				// Buscar en sides del item
				for k := range createRequest.Menu[i].Items[j].Sides {
					if createRequest.Menu[i].Items[j].Sides[k].ID == menuItemID {
						createRequest.Menu[i].Items[j].Sides[k].PhotoUrl = imageURL
						obs.Logger.InfoContext(ctx, "image_url_assigned_to_side", "menuItemId", menuItemID, "imageURL", imageURL)
						found = true
						break
					}
				}
				
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		
		if !found {
			obs.Logger.WarnContext(ctx, "menu_item_not_found_for_image", "menuItemId", menuItemID)
		}
	}
}
