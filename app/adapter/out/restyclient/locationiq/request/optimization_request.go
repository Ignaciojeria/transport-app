package request

import (
	"fmt"
	"strings"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
)

// OptimizationRequest estructura final con los endpoints.
type OptimizationRequest struct {
	ENDPOINTS []Endpoint
}

// Endpoint representa el endpoint y la ruta asociada.
type Endpoint struct {
	Route domain.Route
	URL   string
}

// contains verifica si un slice de strings ya contiene un ítem.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// MapOptimizationRequest construye el request de optimización preservando el orden original.
func MapOptimizationRequest(conf configuration.Conf, plan domain.Plan) OptimizationRequest {

	var endpoints []Endpoint

	// Punto de inicio del plan.
	startPoint := plan.Origin.AddressInfo.Location

	for _, route := range plan.Routes {
		// Punto de destino de la ruta.
		endPoint := route.Destination.AddressInfo.Location

		// Construir slice de coordenadas en orden:
		// 1. Agregar el punto de inicio.
		orderedCoordinates := []string{pointToString(startPoint)}

		// 2. Agregar las coordenadas de las órdenes, en el orden que aparecen.
		for _, order := range route.Orders {
			currentPoint := order.Destination.AddressInfo.Location

			// Excluir si es igual al inicio o al destino.
			if arePointsEqual(currentPoint, startPoint) || arePointsEqual(currentPoint, endPoint) {
				continue
			}

			pointStr := pointToString(currentPoint)
			// Verificar duplicados sin alterar el orden.
			if !contains(orderedCoordinates, pointStr) {
				orderedCoordinates = append(orderedCoordinates, pointStr)
			}
		}

		// 3. Agregar el punto final si es distinto al inicio.
		if !arePointsEqual(startPoint, endPoint) {
			orderedCoordinates = append(orderedCoordinates, pointToString(endPoint))
		}

		// Configurar los parámetros de la API.
		queryParams := []string{
			fmt.Sprintf("key=%s", conf.LOCATION_IQ_ACCESS_TOKEN),
			"source=first",
			"destination=last",
			fmt.Sprintf("roundtrip=%t", arePointsEqual(startPoint, endPoint)),
		}

		// Construir la URL final con las coordenadas en el orden deseado.
		coordinatesString := strings.Join(orderedCoordinates, ";")
		url := fmt.Sprintf("%s/v1/optimize/driving/%s?%s",
			conf.LOCATION_IQ_DNS,
			coordinatesString,
			strings.Join(queryParams, "&"))

		// Crear el endpoint con la ruta completa.
		endpoint := Endpoint{
			Route: route,
			URL:   url,
		}

		endpoints = append(endpoints, endpoint)
	}

	return OptimizationRequest{
		ENDPOINTS: endpoints,
	}
}
