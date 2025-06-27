package vroom

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/vroom/mapper"
	"transport-app/app/adapter/out/vroom/model"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/twpayne/go-polyline"
)

type Optimize func(ctx context.Context, request request.OptimizeFleetRequest) (domain.Plan, error)

func init() {
	ioc.Registry(
		NewOptimize,
		observability.NewObservability,
		NewVroomRestyHeavyClient,
		configuration.NewConf,
	)
}

func NewOptimize(
	obs observability.Observability,
	restyClient *resty.Client,
	conf configuration.Conf,
) Optimize {
	return func(ctx context.Context, req request.OptimizeFleetRequest) (domain.Plan, error) {
		// Convertir el request a la estructura del dominio de optimización
		fleetOptimization := req.Map()

		// Primera optimización - distribución general
		vroomRequest, err := mapper.MapOptimizationRequest(ctx, fleetOptimization)
		if err != nil {
			return domain.Plan{}, err
		}

		obs.Logger.InfoContext(ctx,
			"VROOM_REQUEST",
			"url", conf.VROOM_URL,
			"payload", vroomRequest,
		)

		res, err := restyClient.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetBody(vroomRequest). // Resty hace el marshal automáticamente
			Post(conf.VROOM_URL)

		if err != nil {
			obs.Logger.ErrorContext(ctx,
				"VROOM_REQUEST_ERROR",
				"error", err.Error(),
				"url", conf.VROOM_URL,
			)
			return domain.Plan{}, err
		}

		if res.IsError() {
			obs.Logger.ErrorContext(ctx,
				"VROOM_API_ERROR",
				"status", res.StatusCode(),
				"body", res.String(),
				"request", vroomRequest,
			)

			return domain.Plan{}, fmt.Errorf("VROOM API error (status %d): %s\nRequest payload: %+v",
				res.StatusCode(),
				res.String(),
				vroomRequest)
		}

		// Deserializar la respuesta de VROOM
		var vroomResponse model.VroomOptimizationResponse
		if err := json.Unmarshal(res.Body(), &vroomResponse); err != nil {
			obs.Logger.ErrorContext(ctx,
				"VROOM_RESPONSE_DESERIALIZATION_ERROR",
				"error", err.Error(),
				"body", res.String(),
			)
			return domain.Plan{}, fmt.Errorf("failed to deserialize VROOM response: %w", err)
		}

		plan := domain.Plan{
			ReferenceID: uuid.New().String(),
			PlannedDate: time.Now(),
		}

		// Slice para almacenar todos los polylines consolidados
		var allPolylines []string
		var allRouteData []model.RouteData

		fleetOptimizations, unassignedOrders := vroomResponse.MapOptimizationRequests(ctx, fleetOptimization)

		// Mapear órdenes sin asignar
		for _, unassignedOrder := range unassignedOrders {
			// Mapear las unidades de entrega
			deliveryUnits := make(domain.DeliveryUnits, 0, len(unassignedOrder.DeliveryUnits))
			for _, du := range unassignedOrder.DeliveryUnits {
				// Mapear los items de cada unidad de entrega
				items := make([]domain.Item, 0, len(du.Items))
				for _, item := range du.Items {
					items = append(items, domain.Item{
						Sku: item.Sku,
					})
				}

				// Crear la unidad de entrega del dominio
				deliveryUnit := domain.DeliveryUnit{
					Lpn:   du.Lpn,
					Items: items,
				}
				deliveryUnits = append(deliveryUnits, deliveryUnit)
			}

			// Crear la orden del dominio
			domainOrder := domain.Order{
				ReferenceID:      domain.ReferenceID(unassignedOrder.ReferenceID),
				DeliveryUnits:    deliveryUnits,
				UnassignedReason: "No se pudo asignar a ningún vehículo disponible",
			}

			plan.UnassignedOrders = append(plan.UnassignedOrders, domainOrder)
		}

		// Segunda optimización - optimización individual por ruta
		for optimizationIndex, individualFleetOptimization := range fleetOptimizations {
			individualVroomRequest, err := mapper.MapOptimizationRequest(ctx, individualFleetOptimization)
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Failed to map individual optimization request", "error", err)
				continue
			}

			obs.Logger.InfoContext(ctx,
				"INDIVIDUAL_VROOM_REQUEST",
				"url", conf.VROOM_URL,
				"payload", individualVroomRequest,
			)

			res, err := restyClient.R().
				SetContext(ctx).
				SetHeader("Content-Type", "application/json").
				SetBody(individualVroomRequest).
				Post("http://localhost:3000/")

			if err != nil {
				obs.Logger.ErrorContext(ctx,
					"INDIVIDUAL_VROOM_REQUEST_ERROR",
					"error", err.Error(),
					"url", conf.VROOM_URL,
				)
				continue
			}

			if res.IsError() {
				obs.Logger.ErrorContext(ctx,
					"INDIVIDUAL_VROOM_API_ERROR",
					"status", res.StatusCode(),
					"body", res.String(),
					"request", individualVroomRequest,
				)
				continue
			}

			// Procesar respuesta individual
			var individualVroomResponse model.VroomOptimizationResponse
			if err := json.Unmarshal(res.Body(), &individualVroomResponse); err != nil {
				obs.Logger.ErrorContext(ctx,
					"INDIVIDUAL_VROOM_RESPONSE_DESERIALIZATION_ERROR",
					"error", err.Error(),
					"body", res.String(),
				)
				continue
			}

			// Exportar polyline individual con número secuencial
			polylineFilename := fmt.Sprintf("ui/static/dev/polyline_%03d.json", optimizationIndex+1)
			individualVroomResponse.ExportToPolylineJSON(polylineFilename, fleetOptimization)

			// Consolidar polylines de esta optimización individual
			for _, vroomRoute := range individualVroomResponse.Routes {
				if vroomRoute.Geometry != "" {
					allPolylines = append(allPolylines, vroomRoute.Geometry)
				}

				// También almacenar datos de la ruta para exportación
				routeData := model.RouteData{
					Vehicle:  vroomRoute.Vehicle,
					Cost:     vroomRoute.Cost,
					Duration: vroomRoute.Duration,
					Steps:    make([]model.StepPoint, 0, len(vroomRoute.Steps)),
				}

				// Mapear steps de la ruta
				for stepIndex, step := range vroomRoute.Steps {
					stepPoint := model.StepPoint{
						StepNumber: stepIndex + 1,
						StepType:   step.Type,
						Arrival:    step.Arrival,
					}

					if step.Location != [2]float64{0, 0} {
						stepPoint.Location = step.Location
					}

					if step.Description != "" {
						stepPoint.Description = step.Description
					}

					routeData.Steps = append(routeData.Steps, stepPoint)
				}

				allRouteData = append(allRouteData, routeData)
			}

			obs.Logger.InfoContext(ctx,
				"INDIVIDUAL_OPTIMIZATION_COMPLETED",
				"optimization_index", optimizationIndex+1,
				"routes", len(individualVroomResponse.Routes),
				"unassigned", len(individualVroomResponse.Unassigned),
				"polyline_file", polylineFilename,
			)
		}

		return plan, nil
	}
}

// decodePolyline convierte un string de polyline codificado en coordenadas [lon, lat]
func decodePolyline(polylineStr string) [][]float64 {
	if polylineStr == "" {
		return [][]float64{}
	}

	coords, _, err := polyline.DecodeCoords([]byte(polylineStr))
	if err != nil {
		return [][]float64{}
	}

	// Convertir a formato [lon, lat] que espera Leaflet
	result := make([][]float64, len(coords))
	for i, coord := range coords {
		result[i] = []float64{coord[1], coord[0]} // [lon, lat]
	}

	return result
}
