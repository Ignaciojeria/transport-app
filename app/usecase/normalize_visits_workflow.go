package usecase

import (
	"context"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// NormalizeVisitsWorkflow es un workflow genérico que normaliza las claves de las visitas
// basándose en un mapeo de claves proporcionado
type NormalizeVisitsWorkflow func(ctx context.Context, keyMapping map[string]string, visits []map[string]interface{}) ([]map[string]interface{}, error)

func init() {
	ioc.Registry(
		NewNormalizeVisitsWorkflow,
		observability.NewObservability,
	)
}

func NewNormalizeVisitsWorkflow(
	obs observability.Observability,
) NormalizeVisitsWorkflow {
	return func(ctx context.Context, keyMapping map[string]string, visits []map[string]interface{}) ([]map[string]interface{}, error) {
		obs.Logger.InfoContext(ctx, "Iniciando normalización de visitas", "totalVisits", len(visits), "keyMapping", keyMapping)

		if len(visits) == 0 {
			obs.Logger.InfoContext(ctx, "No hay visitas para normalizar")
			return visits, nil
		}

		// Crear un nuevo array de visitas normalizadas
		normalizedVisits := make([]map[string]interface{}, len(visits))

		// Iterar sobre cada visita original
		for i, originalVisit := range visits {
			obs.Logger.InfoContext(ctx, "Procesando visita", "index", i, "originalVisit", originalVisit)

			// Crear un nuevo mapa para la visita normalizada
			normalizedVisit := make(map[string]interface{})

			// Iterar sobre cada clave en la visita original
			for originalKey, value := range originalVisit {
				if officialKey, exists := keyMapping[originalKey]; exists && officialKey != "" {
					// Usar la clave oficial si existe en el mapeo y no está vacía
					normalizedVisit[officialKey] = value
					obs.Logger.InfoContext(ctx, "Clave renombrada", "originalKey", originalKey, "officialKey", officialKey, "value", value)
				} else {
					// Mantener la clave original si no hay mapeo o está vacía
					normalizedVisit[originalKey] = value
					obs.Logger.InfoContext(ctx, "Clave mantenida", "originalKey", originalKey, "value", value)
				}
			}

			normalizedVisits[i] = normalizedVisit
			obs.Logger.InfoContext(ctx, "Visita normalizada", "index", i, "normalizedVisit", normalizedVisit)
		}

		obs.Logger.InfoContext(ctx, "Visitas normalizadas completadas", "totalVisits", len(normalizedVisits))
		return normalizedVisits, nil
	}
}
