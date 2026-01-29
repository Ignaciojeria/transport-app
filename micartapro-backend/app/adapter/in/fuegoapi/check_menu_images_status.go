package fuegoapi

import (
	"context"
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type ImageStatus struct {
	URL    string `json:"url"`
	Exists bool   `json:"exists"`
}

type MenuImagesStatusResponse struct {
	MenuID        string        `json:"menuId"`
	VersionID    string        `json:"versionId,omitempty"`
	AllReady     bool          `json:"allReady"`     // true si todas las imágenes están disponibles
	PendingCount int           `json:"pendingCount"` // número de imágenes pendientes
	Images       []ImageStatus `json:"images"`      // estado de cada imagen
}

func init() {
	ioc.Registry(
		checkMenuImagesStatus,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetMenuById,
		apimiddleware.NewJWTAuthMiddleware,
		gcs.NewClient,
	)
}

func checkMenuImagesStatus(
	s httpserver.Server,
	obs observability.Observability,
	getMenuById supabaserepo.GetMenuById,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
	gcsClient *storage.Client,
) {
	fuego.Get(s.Manager, "/api/menus/{menuId}/images-status",
		func(c fuego.ContextNoBody) (MenuImagesStatusResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "checkMenuImagesStatus")
			defer span.End()

			menuID := c.PathParam("menuId")
			if menuID == "" {
				return MenuImagesStatusResponse{}, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			versionID := c.QueryParam("version_id")

			// Obtener el menú
			menu, err := getMenuById(spanCtx, menuID, versionID)
			if err != nil {
				if err == supabaserepo.ErrMenuNotFound {
					return MenuImagesStatusResponse{}, fuego.HTTPError{
						Title:  "menu not found",
						Detail: "menu with the provided menu_id was not found",
						Status: http.StatusNotFound,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error getting menu", "error", err)
				return MenuImagesStatusResponse{}, fuego.HTTPError{
					Title:  "error getting menu",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			// Normalizar URLs de imágenes por si se guardaron mal en BD (httpshttps, https.storage, etc.)
			normalizeMenuImageURLs(&menu)

			// Extraer todas las URLs de imágenes del menú (ya normalizadas)
			imageURLs := extractImageURLs(menu)

			// Verificar existencia de cada imagen en GCS
			imageStatuses := make([]ImageStatus, 0, len(imageURLs))
			pendingCount := 0

			for _, url := range imageURLs {
				if url == "" {
					continue
				}

				// Las URLs ya están normalizadas por normalizeMenuImageURLs, pero aplicamos normalizeGCSURL
				// por si acaso (es idempotente, no causa problemas)
				normalizedURL := normalizeGCSURL(url)

				exists := false
				// Solo verificar URLs de GCS
				if strings.Contains(normalizedURL, "storage.googleapis.com") {
					exists = checkImageExists(spanCtx, gcsClient, normalizedURL, obs)
				} else {
					// Para URLs no-GCS, asumir que existen (no podemos verificarlas)
					exists = true
				}

				imageStatuses = append(imageStatuses, ImageStatus{
					URL:    normalizedURL,
					Exists: exists,
				})

				if !exists {
					pendingCount++
				}
			}

			allReady := pendingCount == 0

			obs.Logger.InfoContext(spanCtx, "menu_images_status_checked",
				"menuID", menuID,
				"versionID", versionID,
				"totalImages", len(imageURLs),
				"pendingCount", pendingCount,
				"allReady", allReady)

			return MenuImagesStatusResponse{
				MenuID:       menuID,
				VersionID:   versionID,
				AllReady:    allReady,
				PendingCount: pendingCount,
				Images:      imageStatuses,
			}, nil
		},
		option.Summary("Check menu images status"),
		option.Description("Verifies which images from a menu exist in GCS"),
		option.Tags("menu", "images"),
		option.Middleware(jwtAuthMiddleware),
	)
}

// extractImageURLs extrae todas las URLs de imágenes del menú
func extractImageURLs(menu events.MenuCreateRequest) []string {
	urls := make([]string, 0)

	// Cover image
	if menu.CoverImage != "" {
		urls = append(urls, menu.CoverImage)
	}

	// Footer image
	if menu.FooterImage != "" {
		urls = append(urls, menu.FooterImage)
	}

	// Images de items y sides
	for _, category := range menu.Menu {
		for _, item := range category.Items {
			if item.PhotoUrl != "" {
				urls = append(urls, item.PhotoUrl)
			}
			for _, side := range item.Sides {
				if side.PhotoUrl != "" {
					urls = append(urls, side.PhotoUrl)
				}
			}
		}
	}

	return urls
}

// normalizeGCSURL corrige URLs de GCS mal formateadas
// Es idempotente: puede aplicarse múltiples veces sin causar efectos secundarios
func normalizeGCSURL(url string) string {
	if url == "" {
		return url
	}
	
	// Verificar si la URL ya está correctamente formateada
	if strings.HasPrefix(url, "https://storage.googleapis.com") || strings.HasPrefix(url, "http://storage.googleapis.com") {
		return url
	}
	
	// Corregir "https.storage.googleapis.com" → "https://storage.googleapis.com"
	// Solo aplicar si el patrón incorrecto está presente
	if strings.Contains(url, "https.storage.googleapis.com") {
		url = strings.ReplaceAll(url, "https.storage.googleapis.com", "https://storage.googleapis.com")
	}
	if strings.Contains(url, "http.storage.googleapis.com") {
		url = strings.ReplaceAll(url, "http.storage.googleapis.com", "http://storage.googleapis.com")
	}
	
	// Manejar casos donde se duplicó "https" (httpshttps://storage...)
	url = strings.ReplaceAll(url, "httpshttps://", "https://")
	url = strings.ReplaceAll(url, "httphttp://", "http://")
	
	return url
}

// checkImageExists verifica si una imagen existe en GCS haciendo un HEAD request
func checkImageExists(ctx context.Context, client *storage.Client, imageURL string, obs observability.Observability) bool {
	if client == nil {
		return false
	}

	// Extraer bucket y object path de la URL
	// Formato: https://storage.googleapis.com/<bucket>/<object-path>
	if !strings.Contains(imageURL, "storage.googleapis.com/") {
		return false
	}

	parts := strings.Split(imageURL, "storage.googleapis.com/")
	if len(parts) != 2 {
		return false
	}

	pathParts := strings.SplitN(parts[1], "/", 2)
	if len(pathParts) != 2 {
		return false
	}

	bucketName := pathParts[0]
	objectPath := pathParts[1]

	// Verificar si el objeto existe usando Attrs
	_, err := client.Bucket(bucketName).Object(objectPath).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false
		}
		obs.Logger.WarnContext(ctx, "error_checking_image_exists", "error", err, "url", imageURL)
		return false
	}

	return true
}
