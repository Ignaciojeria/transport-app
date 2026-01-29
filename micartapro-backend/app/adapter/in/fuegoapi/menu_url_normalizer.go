package fuegoapi

import "micartapro/app/events"

// normalizeMenuImageURLs normaliza todas las URLs de imágenes del menú (CoverImage, FooterImage, PhotoUrl en items y sides).
// Corrige formatos mal guardados en BD como "httpshttps://" o "https.storage.googleapis.com".
// Es idempotente y seguro llamarlo múltiples veces.
// Usa el método NormalizeImageURLs del tipo MenuCreateRequest.
func normalizeMenuImageURLs(menu *events.MenuCreateRequest) {
	if menu == nil {
		return
	}
	menu.NormalizeImageURLs()
}
