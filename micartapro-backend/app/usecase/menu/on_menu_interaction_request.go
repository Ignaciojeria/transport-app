package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/storage"
	"micartapro/app/adapter/out/agents"
	"micartapro/app/adapter/out/imagegenerator"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"strings"

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
		imagegenerator.NewImageEditor,
		gcs.NewClient)
}

func NewOnMenuInteractionRequest(
	obs observability.Observability,
	menuInteractionAgent agents.MenuInteractionAgent,
	publisherManager eventprocessing.PublisherManager,
	getMenuById supabaserepo.GetMenuById,
	generateImage imagegenerator.GenerateImage,
	editImage imagegenerator.EditImage,
	gcsClient *storage.Client) OnMenuInteractionRequest {
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

			// Pre-firmar URLs y asignar placeholders para imágenes que se generarán de forma asíncrona
			// Esto permite guardar el menú rápidamente sin esperar la generación de imágenes
			if err := prepareImagePlaceholders(ctx, obs, gcsClient, &createRequest, menu); err != nil {
				obs.Logger.ErrorContext(ctx, "error_preparing_image_placeholders", "error", err)
				return "", fmt.Errorf("error preparando placeholders de imágenes: %w", err)
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

// prepareImagePlaceholders pre-firma URLs y asigna placeholders para imágenes que se generarán de forma asíncrona
func prepareImagePlaceholders(
	ctx context.Context,
	obs observability.Observability,
	gcsClient *storage.Client,
	createRequest *events.MenuCreateRequest,
	menu events.MenuCreateRequest,
) error {
	spanCtx, span := obs.Tracer.Start(ctx, "prepare_image_placeholders")
	defer span.End()

	userID, ok := sharedcontext.UserIDFromContext(spanCtx)
	if !ok || userID == "" {
		return fmt.Errorf("userID is required but not found in context")
	}

	// Procesar imagen de portada si hay solicitud de generación
	if createRequest.CoverImageGenerationRequest != nil {
		fileName := fmt.Sprintf("cover-%s.png", uuid.New().String()[:8])
		uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, "image/png")
		if err != nil {
			return fmt.Errorf("error generando signed URL para cover: %w", err)
		}
		// Asignar placeholder (la URL pública donde se guardará la imagen)
		// Normalizar URL antes de asignar
		createRequest.CoverImage = events.NormalizeGCSURL(publicURL)
		// Guardar URLs en campos temporales para publicar evento después
		createRequest.CoverImageGenerationRequest.UploadURL = uploadURL
		createRequest.CoverImageGenerationRequest.PublicURL = publicURL
		obs.Logger.InfoContext(spanCtx, "cover_image_placeholder_prepared", "publicURL", publicURL)
	}

	// Procesar imagen de portada si hay solicitud de edición
	if createRequest.CoverImageEditionRequest != nil {
		// Si no hay referenceImageUrl, usar la coverImage del menú actual
		if createRequest.CoverImageEditionRequest.ReferenceImageUrl == "" && menu.CoverImage != "" {
			createRequest.CoverImageEditionRequest.ReferenceImageUrl = menu.CoverImage
		}

		fileName := fmt.Sprintf("cover-edited-%s.png", uuid.New().String()[:8])
		contentType := "image/png"
		if strings.Contains(strings.ToLower(fileName), ".jpg") || strings.Contains(strings.ToLower(fileName), ".jpeg") {
			contentType = "image/jpeg"
		}

		uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, contentType)
		if err != nil {
			return fmt.Errorf("error generando signed URL para cover editado: %w", err)
		}
		// Asignar placeholder
		// Normalizar URL antes de asignar
		createRequest.CoverImage = events.NormalizeGCSURL(publicURL)
		// Guardar URLs en campos temporales
		createRequest.CoverImageEditionRequest.UploadURL = uploadURL
		createRequest.CoverImageEditionRequest.PublicURL = publicURL
		obs.Logger.InfoContext(spanCtx, "cover_image_edition_placeholder_prepared", "publicURL", publicURL)
	}

	// Procesar imágenes de items si hay solicitudes de generación
	for i := range createRequest.ImageGenerationRequests {
		req := &createRequest.ImageGenerationRequests[i]
		fileName := fmt.Sprintf("%s-%s.png", req.MenuItemID, uuid.New().String()[:8])
		uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, "image/png")
		if err != nil {
			return fmt.Errorf("error generando signed URL para item %s: %w", req.MenuItemID, err)
		}
		// Guardar URLs en campos temporales
		req.UploadURL = uploadURL
		req.PublicURL = publicURL
		// Asignar placeholder al item del menú (para que se persista)
		assignImageURLToMenuItem(createRequest, req.MenuItemID, publicURL, obs, spanCtx)
		obs.Logger.InfoContext(spanCtx, "item_image_placeholder_prepared", "menuItemId", req.MenuItemID, "publicURL", publicURL)
	}

	// Procesar imágenes de items si hay solicitudes de edición
	for i := range createRequest.ImageEditionRequests {
		req := &createRequest.ImageEditionRequests[i]
		// Si no hay referenceImageUrl, buscar la imagen del item en el menú actual
		if req.ReferenceImageUrl == "" {
			imageURL := findImageUrlInMenu(menu, req.MenuItemID)
			if imageURL != "" {
				req.ReferenceImageUrl = imageURL
			}
		}

		fileName := fmt.Sprintf("%s-edited-%s.png", req.MenuItemID, uuid.New().String()[:8])
		contentType := "image/png"
		if strings.Contains(strings.ToLower(fileName), ".jpg") || strings.Contains(strings.ToLower(fileName), ".jpeg") {
			contentType = "image/jpeg"
		}

		uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, contentType)
		if err != nil {
			return fmt.Errorf("error generando signed URL para item editado %s: %w", req.MenuItemID, err)
		}
		// Guardar URLs en campos temporales
		req.UploadURL = uploadURL
		req.PublicURL = publicURL
		// Asignar placeholder al item del menú (para que se persista)
		assignImageURLToMenuItem(createRequest, req.MenuItemID, publicURL, obs, spanCtx)
		obs.Logger.InfoContext(spanCtx, "item_image_edition_placeholder_prepared", "menuItemId", req.MenuItemID, "publicURL", publicURL)
	}

	return nil
}

// assignImageURLToMenuItem asigna una URL de imagen a un item o side del menú
// Normaliza la URL antes de asignarla para evitar formatos incorrectos
func assignImageURLToMenuItem(
	createRequest *events.MenuCreateRequest,
	menuItemID string,
	imageURL string,
	obs observability.Observability,
	ctx context.Context,
) {
	// Normalizar URL antes de asignar
	normalizedURL := events.NormalizeGCSURL(imageURL)
	
	// Manejar IDs especiales
	if menuItemID == "footer" {
		createRequest.FooterImage = normalizedURL
		obs.Logger.InfoContext(ctx, "image_url_assigned_to_footer", "imageURL", normalizedURL)
		return
	}

	// Buscar el item o side en el menú
	found := false
	for i := range createRequest.Menu {
		for j := range createRequest.Menu[i].Items {
			if createRequest.Menu[i].Items[j].ID == menuItemID {
				createRequest.Menu[i].Items[j].PhotoUrl = normalizedURL
				obs.Logger.InfoContext(ctx, "image_url_assigned_to_item", "menuItemId", menuItemID, "imageURL", normalizedURL)
				found = true
				break
			}

			// Buscar en sides del item
			for k := range createRequest.Menu[i].Items[j].Sides {
				if createRequest.Menu[i].Items[j].Sides[k].ID == menuItemID {
					createRequest.Menu[i].Items[j].Sides[k].PhotoUrl = normalizedURL
					obs.Logger.InfoContext(ctx, "image_url_assigned_to_side", "menuItemId", menuItemID, "imageURL", normalizedURL)
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

// processCoverImageGenerationRequest procesa la solicitud de generación de imagen de portada
func processCoverImageGenerationRequest(
	ctx context.Context,
	obs observability.Observability,
	gcsClient *storage.Client,
	createRequest *events.MenuCreateRequest,
	generateImage imagegenerator.GenerateImage,
) error {
	spanCtx, span := obs.Tracer.Start(ctx, "process_cover_image")
	defer span.End()

	coverReq := createRequest.CoverImageGenerationRequest
	obs.Logger.InfoContext(spanCtx, "generating_cover_image", "prompt", coverReq.Prompt)

	userID, ok := sharedcontext.UserIDFromContext(spanCtx)
	if !ok || userID == "" {
		return fmt.Errorf("userID is required but not found in context")
	}

	// Crear nombre de archivo único
	fileName := fmt.Sprintf("cover-%s.png", uuid.New().String()[:8])

	// Generar URLs pre-firmadas
	uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, "image/png")
	if err != nil {
		obs.Logger.ErrorContext(spanCtx, "error_generating_signed_url", "error", err)
		return fmt.Errorf("error generando signed URL: %w", err)
	}

	// Generar la imagen y subir usando la signed URL (retorna URL pública)
	publicURL, err = generateImage(spanCtx, coverReq.Prompt, "16:9", coverReq.ImageCount, uploadURL, publicURL)
	if err != nil {
		obs.Logger.ErrorContext(spanCtx, "error_generating_cover_image", "error", err)
		return fmt.Errorf("error generando imagen de portada: %w", err)
	}

	// Asignar la URL directamente al coverImage
	// Normalizar URL antes de asignar
	createRequest.CoverImage = events.NormalizeGCSURL(publicURL)
	obs.Logger.InfoContext(spanCtx, "cover_image_processed_successfully", "publicURL", createRequest.CoverImage)

	return nil
}

// processImageGenerationRequests procesa todas las solicitudes de generación de imágenes en paralelo
func processImageGenerationRequests(
	ctx context.Context,
	obs observability.Observability,
	gcsClient *storage.Client,
	createRequest *events.MenuCreateRequest,
	generateImage imagegenerator.GenerateImage,
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

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				errChan <- fmt.Errorf("userID is required but not found in context")
				return
			}

			// Crear nombre de archivo único
			fileName := fmt.Sprintf("%s-%s.png", req.MenuItemID, uuid.New().String()[:8])

			// Generar URLs pre-firmadas
			uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, "image/png")
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_generating_signed_url", "error", err, "menuItemId", req.MenuItemID)
				errChan <- fmt.Errorf("error generando signed URL para %s: %w", req.MenuItemID, err)
				return
			}

			// Generar la imagen y subir usando la signed URL (retorna URL pública)
			publicURL, err = generateImage(spanCtx, req.Prompt, req.AspectRatio, req.ImageCount, uploadURL, publicURL)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_generating_image", "error", err, "menuItemId", req.MenuItemID)
				errChan <- fmt.Errorf("error generando imagen para %s: %w", req.MenuItemID, err)
				return
			}

			// Guardar la URL en el mapa
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

// findImageUrlInMenu busca la URL de la imagen de un item o side en el menú actual
func findImageUrlInMenu(menu events.MenuCreateRequest, menuItemID string) string {
	// Buscar en items
	for _, category := range menu.Menu {
		for _, item := range category.Items {
			if item.ID == menuItemID {
				if item.PhotoUrl != "" {
					return item.PhotoUrl
				}
			}

			// Buscar en sides del item
			for _, side := range item.Sides {
				if side.ID == menuItemID {
					if side.PhotoUrl != "" {
						return side.PhotoUrl
					}
				}
			}
		}
	}

	// Manejar IDs especiales
	if menuItemID == "footer" && menu.FooterImage != "" {
		return menu.FooterImage
	}

	return ""
}

// updateMenuWithImageURLs actualiza el menú asignando las URLs de las imágenes a los items/sides correspondientes
// También maneja IDs especiales: "footer" para footerImage (cover ahora se maneja por separado)
// Normaliza las URLs antes de asignarlas para evitar formatos incorrectos
func updateMenuWithImageURLs(
	createRequest *events.MenuCreateRequest,
	imageURLs map[string]string,
	obs observability.Observability,
	ctx context.Context,
) {
	for menuItemID, imageURL := range imageURLs {
		// Normalizar URL antes de asignar
		normalizedURL := events.NormalizeGCSURL(imageURL)
		
		// Manejar IDs especiales para imágenes del menú
		if menuItemID == "footer" {
			createRequest.FooterImage = normalizedURL
			obs.Logger.InfoContext(ctx, "image_url_assigned_to_footer", "imageURL", normalizedURL)
			continue
		}

		// Buscar el item o side en el menú
		found := false

		// Buscar en items
		for i := range createRequest.Menu {
			for j := range createRequest.Menu[i].Items {
				if createRequest.Menu[i].Items[j].ID == menuItemID {
					createRequest.Menu[i].Items[j].PhotoUrl = normalizedURL
					obs.Logger.InfoContext(ctx, "image_url_assigned_to_item", "menuItemId", menuItemID, "imageURL", normalizedURL)
					found = true
					break
				}

				// Buscar en sides del item
				for k := range createRequest.Menu[i].Items[j].Sides {
					if createRequest.Menu[i].Items[j].Sides[k].ID == menuItemID {
						createRequest.Menu[i].Items[j].Sides[k].PhotoUrl = normalizedURL
						obs.Logger.InfoContext(ctx, "image_url_assigned_to_side", "menuItemId", menuItemID, "imageURL", normalizedURL)
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

// processCoverImageEditionRequest procesa la solicitud de edición de imagen de portada
func processCoverImageEditionRequest(
	ctx context.Context,
	obs observability.Observability,
	gcsClient *storage.Client,
	createRequest *events.MenuCreateRequest,
	editImage imagegenerator.EditImage,
) error {
	spanCtx, span := obs.Tracer.Start(ctx, "process_cover_image_edition")
	defer span.End()

	coverReq := createRequest.CoverImageEditionRequest
	obs.Logger.InfoContext(spanCtx, "editing_cover_image", "prompt", coverReq.Prompt, "referenceImageUrl", coverReq.ReferenceImageUrl)

	userID, ok := sharedcontext.UserIDFromContext(spanCtx)
	if !ok || userID == "" {
		return fmt.Errorf("userID is required but not found in context")
	}

	// Crear nombre de archivo único
	fileName := fmt.Sprintf("cover-edited-%s.png", uuid.New().String()[:8])

	// Determinar contentType basado en la extensión
	contentType := "image/png"
	if strings.Contains(strings.ToLower(fileName), ".jpg") || strings.Contains(strings.ToLower(fileName), ".jpeg") {
		contentType = "image/jpeg"
	}

	// Generar URLs pre-firmadas
	uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, contentType)
	if err != nil {
		obs.Logger.ErrorContext(spanCtx, "error_generating_signed_url", "error", err)
		return fmt.Errorf("error generando signed URL: %w", err)
	}

	// Editar la imagen y subir usando la signed URL (retorna URL pública)
	// Pasar "cover" como menuItemId para indicar que es la imagen de portada
	publicURL, err = editImage(spanCtx, coverReq.Prompt, coverReq.ReferenceImageUrl, "16:9", coverReq.ImageCount, "cover", uploadURL, publicURL)
	if err != nil {
		obs.Logger.ErrorContext(spanCtx, "error_editing_cover_image", "error", err)
		return fmt.Errorf("error editando imagen de portada: %w", err)
	}

	// Asignar la URL directamente al coverImage (nueva versión editada)
	// Normalizar URL antes de asignar
	createRequest.CoverImage = events.NormalizeGCSURL(publicURL)
	obs.Logger.InfoContext(spanCtx, "cover_image_edited_successfully", "publicURL", createRequest.CoverImage)

	return nil
}

// processImageEditionRequests procesa todas las solicitudes de edición de imágenes en paralelo
func processImageEditionRequests(
	ctx context.Context,
	obs observability.Observability,
	gcsClient *storage.Client,
	createRequest *events.MenuCreateRequest,
	editImage imagegenerator.EditImage,
) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(createRequest.ImageEditionRequests))

	// Mapa para almacenar las URLs editadas por menuItemId
	imageURLs := make(map[string]string)
	var mu sync.Mutex

	for _, imgReq := range createRequest.ImageEditionRequests {
		wg.Add(1)
		go func(req events.ImageEditionRequest) {
			defer wg.Done()

			spanCtx, span := obs.Tracer.Start(ctx, "process_single_image_edition")
			defer span.End()

			obs.Logger.InfoContext(spanCtx, "editing_image_for_item", "menuItemId", req.MenuItemID, "prompt", req.Prompt, "referenceImageUrl", req.ReferenceImageUrl)

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				errChan <- fmt.Errorf("userID is required but not found in context")
				return
			}

			// Crear nombre de archivo único para la nueva versión editada
			fileName := fmt.Sprintf("%s-edited-%s.png", req.MenuItemID, uuid.New().String()[:8])

			// Determinar contentType basado en la extensión
			contentType := "image/png"
			if strings.Contains(strings.ToLower(fileName), ".jpg") || strings.Contains(strings.ToLower(fileName), ".jpeg") {
				contentType = "image/jpeg"
			}

			// Generar URLs pre-firmadas
			uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURL(spanCtx, gcsClient, obs, userID, fileName, contentType)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_generating_signed_url", "error", err, "menuItemId", req.MenuItemID)
				errChan <- fmt.Errorf("error generando signed URL para %s: %w", req.MenuItemID, err)
				return
			}

			// Editar la imagen y subir usando la signed URL (retorna URL pública)
			publicURL, err = editImage(spanCtx, req.Prompt, req.ReferenceImageUrl, req.AspectRatio, req.ImageCount, req.MenuItemID, uploadURL, publicURL)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_editing_image", "error", err, "menuItemId", req.MenuItemID)
				errChan <- fmt.Errorf("error editando imagen para %s: %w", req.MenuItemID, err)
				return
			}

			// Guardar la URL en el mapa
			mu.Lock()
			imageURLs[req.MenuItemID] = publicURL
			mu.Unlock()

			obs.Logger.InfoContext(spanCtx, "image_edited_successfully", "menuItemId", req.MenuItemID, "publicURL", publicURL)
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

	// 5. Actualizar el menú con las URLs de las imágenes editadas
	updateMenuWithImageURLs(createRequest, imageURLs, obs, ctx)

	return nil
}
