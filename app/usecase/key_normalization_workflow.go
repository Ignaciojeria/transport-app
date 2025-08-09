package usecase

import (
	"context"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// KeyNormalizationWorkflow es un workflow genérico que normaliza las claves de cualquier estructura
// basándose en un mapeo de claves proporcionado
type KeyNormalizationWorkflow func(ctx context.Context, keyMapping map[string]string, data []map[string]interface{}) ([]map[string]interface{}, error)

func init() {
	ioc.Registry(
		NewKeyNormalizationWorkflow,
		observability.NewObservability,
	)
}

func NewKeyNormalizationWorkflow(
	obs observability.Observability,
) KeyNormalizationWorkflow {
	return func(ctx context.Context, keyMapping map[string]string, data []map[string]interface{}) ([]map[string]interface{}, error) {
		obs.Logger.InfoContext(ctx, "Iniciando normalización de claves", "totalItems", len(data), "keyMapping", keyMapping)

		if len(data) == 0 {
			obs.Logger.InfoContext(ctx, "No hay elementos para normalizar")
			return data, nil
		}

		// Invertir el mapeo: el LLM devuelve officialKey -> originalKey, pero necesitamos originalKey -> officialKey
		invertedMapping := make(map[string]string)
		for officialKey, originalKey := range keyMapping {
			if originalKey != "" && officialKey != "" {
				invertedMapping[originalKey] = officialKey
			}
		}

		obs.Logger.InfoContext(ctx, "Mapeo invertido", "invertedMapping", invertedMapping)

		// Crear un nuevo array de elementos normalizados
		normalizedData := make([]map[string]interface{}, len(data))

		// Iterar sobre cada elemento original
		for i, originalItem := range data {
			obs.Logger.InfoContext(ctx, "Procesando elemento", "index", i, "originalItem", originalItem)

			// Crear un nuevo mapa para el elemento normalizado
			normalizedItem := make(map[string]interface{})

			// Iterar sobre cada clave en el elemento original
			for originalKey, value := range originalItem {
				if officialKey, exists := invertedMapping[originalKey]; exists && officialKey != "" {
					// Usar la clave oficial si existe en el mapeo y no está vacía
					normalizedItem[officialKey] = value
					obs.Logger.InfoContext(ctx, "Clave renombrada", "originalKey", originalKey, "officialKey", officialKey, "value", value)
				} else {
					// Mantener la clave original si no hay mapeo o está vacía
					normalizedItem[originalKey] = value
					obs.Logger.InfoContext(ctx, "Clave mantenida", "originalKey", originalKey, "value", value)
				}
			}

			normalizedData[i] = normalizedItem
			obs.Logger.InfoContext(ctx, "Elemento normalizado", "index", i, "normalizedItem", normalizedItem)
		}

		obs.Logger.InfoContext(ctx, "Normalización completada", "totalItems", len(normalizedData))
		return normalizedData, nil
	}
}
