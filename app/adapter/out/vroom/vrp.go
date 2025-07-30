package vroom

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/out/vroom/mapper"
	"transport-app/app/adapter/out/vroom/model"
	"transport-app/app/domain"
	"transport-app/app/domain/optimization"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/twpayne/go-polyline"
)

type Optimize func(ctx context.Context, request optimization.FleetOptimization) (domain.Plan, error)

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
	return func(ctx context.Context, fleetOptimization optimization.FleetOptimization) (domain.Plan, error) {

		vroomRequest, err := mapper.MapOptimizationRequest(ctx, fleetOptimization)
		if err != nil {
			return domain.Plan{}, err
		}
		/*
			obs.Logger.InfoContext(ctx,
				"VROOM_REQUEST",
				"url", conf.VROOM_PLANNER_URL,
				"payload", vroomRequest,
			)*/

		res, err := restyClient.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetBody(vroomRequest). // Resty hace el marshal automáticamente
			Post(conf.VROOM_PLANNER_URL)

		if err != nil {
			obs.Logger.ErrorContext(ctx,
				"VROOM_REQUEST_ERROR",
				"error", err.Error(),
				"url", conf.VROOM_PLANNER_URL,
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

		// Mapear órdenes sin asignar usando el método del dominio
		var domainUnassignedOrders []domain.Order
		for _, unassignedOrder := range unassignedOrders {
			domainOrder := createOrderFromOptimizationOrder(unassignedOrder)
			domainUnassignedOrders = append(domainUnassignedOrders, domainOrder)
		}

		// Usar el método del dominio para agregar órdenes sin asignar
		plan.AddUnassignedOrders(domainUnassignedOrders, "No se pudo asignar a ningún vehículo disponible")

		// Segunda optimización - optimización individual por ruta
		for optimizationIndex, individualFleetOptimization := range fleetOptimizations {
			individualVroomRequest, err := mapper.MapOptimizationRequest(ctx, individualFleetOptimization)
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Failed to map individual optimization request", "error", err)
				continue
			}
			/*
				obs.Logger.InfoContext(ctx,
					"INDIVIDUAL_VROOM_REQUEST",
					"url", conf.VROOM_OPTIMIZER_URL,
					"payload", individualVroomRequest,
				)*/

			res, err := restyClient.R().
				SetContext(ctx).
				SetHeader("Content-Type", "application/json").
				SetBody(individualVroomRequest).
				Post(conf.VROOM_OPTIMIZER_URL)

			if err != nil {
				obs.Logger.ErrorContext(ctx,
					"INDIVIDUAL_VROOM_REQUEST_ERROR",
					"error", err.Error(),
					"url", conf.VROOM_OPTIMIZER_URL,
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

			// Mapear cada ruta optimizada al plan
			for _, vroomRoute := range individualVroomResponse.Routes {
				// Crear la ruta del dominio con todos los datos
				route := domain.Route{
					ReferenceID: uuid.New().String(),
				}

				// Mapear vehículo con información completa
				if len(individualFleetOptimization.Vehicles) > 0 {
					vehicle := individualFleetOptimization.Vehicles[0] // Solo hay un vehículo por optimización individual
					route.Vehicle = domain.Vehicle{
						Plate: vehicle.Plate,
						Weight: struct {
							Value         int
							UnitOfMeasure string
						}{
							Value:         int(vehicle.Capacity.Weight),
							UnitOfMeasure: "g",
						},
						Insurance: struct {
							PolicyStartDate      string
							PolicyExpirationDate string
							PolicyRenewalDate    string
							MaxInsuranceCoverage struct {
								Amount   float64
								Currency string
							}
						}{
							MaxInsuranceCoverage: struct {
								Amount   float64
								Currency string
							}{
								Amount:   float64(vehicle.Capacity.Insurance),
								Currency: "CLP",
							},
						},
					}
				}

				// Mapear origen y destino basado en los steps con información completa
				if len(vroomRoute.Steps) > 0 {
					// Origen: primer step con location
					for _, step := range vroomRoute.Steps {
						if step.Type == "start" && len(step.Location) == 2 {
							route.Origin = domain.NodeInfo{
								ReferenceID: domain.ReferenceID(uuid.New().String()),
								AddressInfo: domain.AddressInfo{
									Coordinates: domain.Coordinates{
										Point: orb.Point{step.Location[0], step.Location[1]},
									},
								},
							}
							break
						}
					}

					// Destino: último step con location
					for i := len(vroomRoute.Steps) - 1; i >= 0; i-- {
						step := vroomRoute.Steps[i]
						if step.Type == "end" && len(step.Location) == 2 {
							route.Destination = domain.NodeInfo{
								ReferenceID: domain.ReferenceID(uuid.New().String()),
								AddressInfo: domain.AddressInfo{
									Coordinates: domain.Coordinates{
										Point: orb.Point{step.Location[0], step.Location[1]},
									},
								},
							}
							break
						}
					}
				}

				// Mapear órdenes basadas en los jobs y shipments de los steps
				var routeOrders []domain.Order
				for _, step := range vroomRoute.Steps {
					// Manejar jobs (solo delivery)
					if step.Job > 0 {
						visit := findVisitByJobID(step.Job, individualFleetOptimization.Visits)
						if visit != nil {
							orders := createOrdersFromVisit(visit, false)
							routeOrders = append(routeOrders, orders...)
						}
					}

					// Manejar shipments (pickup y delivery)
					if step.Shipment > 0 {
						visit := findVisitByShipmentID(step.Shipment, individualFleetOptimization.Visits)
						if visit != nil {
							orders := createOrdersFromVisit(visit, true)
							routeOrders = append(routeOrders, orders...)
						}
					}
				}
				route.Orders = routeOrders

				// Mapear TimeWindow si está disponible en el vehículo
				if len(individualFleetOptimization.Vehicles) > 0 {
					vehicle := individualFleetOptimization.Vehicles[0]
					route.TimeWindow = domain.TimeWindow{
						Start: vehicle.TimeWindow.Start,
						End:   vehicle.TimeWindow.End,
					}
				}

				// Mapear geometría de la ruta
				if vroomRoute.Geometry != "" {
					route.Geometry = domain.RouteGeometry{
						Encoding: "polyline",
						Type:     "linestring",
						Value:    vroomRoute.Geometry,
					}
				}

				// Agregar la ruta al plan
				plan.Routes = append(plan.Routes, route)

				// Consolidar polylines de esta optimización individual
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
			/*
				obs.Logger.InfoContext(ctx,
					"INDIVIDUAL_OPTIMIZATION_COMPLETED",
					"optimization_index", optimizationIndex+1,
					"routes", len(individualVroomResponse.Routes),
					"unassigned", len(individualVroomResponse.Unassigned),
					"polyline_file", polylineFilename,
				)
			*/
		}

		return plan, nil
	}
}

// findVisitByJobID busca una visita que corresponde a un job (solo delivery válido)
func findVisitByJobID(jobID int64, visits []optimization.Visit) *optimization.Visit {
	// Los job IDs en VROOM corresponden al índice de la visita en el request original
	// jobID 1 = primera visita, jobID 2 = segunda visita, etc.
	// Pero necesitamos ajustar porque los jobs se crean solo para visitas sin pickup válido

	jobIndex := 0
	for i, v := range visits {
		// Verificar si esta visita tiene solo delivery válido (job)
		hasValidPickup := v.Pickup.AddressInfo.Coordinates.Longitude != 0 || v.Pickup.AddressInfo.Coordinates.Latitude != 0
		hasValidDelivery := v.Delivery.AddressInfo.Coordinates.Longitude != 0 || v.Delivery.AddressInfo.Coordinates.Latitude != 0

		if !hasValidPickup && hasValidDelivery {
			// Esta visita corresponde a un job
			jobIndex++
			if int64(jobIndex) == jobID {
				return &visits[i]
			}
		}
	}
	return nil
}

// findVisitByShipmentID busca una visita que corresponde a un shipment (pickup y delivery válidos)
func findVisitByShipmentID(shipmentID int64, visits []optimization.Visit) *optimization.Visit {
	// Los shipment IDs en VROOM corresponden al índice de la visita en el request original
	// shipmentID 1 = primera visita, shipmentID 2 = segunda visita, etc.
	// Pero necesitamos ajustar porque los shipments se crean solo para visitas con pickup válido

	shipmentIndex := 0
	for i, v := range visits {
		// Verificar si esta visita tiene pickup y delivery válidos (shipment)
		hasValidPickup := v.Pickup.AddressInfo.Coordinates.Longitude != 0 || v.Pickup.AddressInfo.Coordinates.Latitude != 0
		hasValidDelivery := v.Delivery.AddressInfo.Coordinates.Longitude != 0 || v.Delivery.AddressInfo.Coordinates.Latitude != 0

		if hasValidPickup && hasValidDelivery {
			// Esta visita corresponde a un shipment
			shipmentIndex++
			if int64(shipmentIndex) == shipmentID {
				return &visits[i]
			}
		}
	}
	return nil
}

// createOrdersFromVisit crea órdenes del dominio basadas en una visita
func createOrdersFromVisit(visit *optimization.Visit, hasPickup bool) []domain.Order {
	var orders []domain.Order

	// Crear orden para cada order en la visita
	for _, orderReq := range visit.Orders {
		order := domain.Order{
			ReferenceID: domain.ReferenceID(orderReq.ReferenceID),
		}

		// Mapear destino con información completa
		order.Destination = domain.NodeInfo{
			ReferenceID: domain.ReferenceID(uuid.New().String()),
			AddressInfo: domain.AddressInfo{
				// Información de contacto del destino
				Contact: domain.Contact{
					FullName:     visit.Delivery.AddressInfo.Contact.FullName,
					PrimaryEmail: visit.Delivery.AddressInfo.Contact.Email,
					PrimaryPhone: visit.Delivery.AddressInfo.Contact.Phone,
					NationalID:   visit.Delivery.AddressInfo.Contact.NationalID,
				},
				// Información política/geográfica
				PoliticalArea: domain.PoliticalArea{
					Code:            visit.Delivery.AddressInfo.PoliticalArea.Code,
					AdminAreaLevel1: visit.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel1,
					AdminAreaLevel2: visit.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel2,
					AdminAreaLevel3: visit.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel3,
					AdminAreaLevel4: visit.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel4,
				},
				// Información de dirección
				AddressLine1: visit.Delivery.AddressInfo.AddressLine1,
				AddressLine2: visit.Delivery.AddressInfo.AddressLine2,
				ZipCode:      visit.Delivery.AddressInfo.ZipCode,
				// Coordenadas
				Coordinates: domain.Coordinates{
					Point: orb.Point{
						visit.Delivery.AddressInfo.Coordinates.Longitude,
						visit.Delivery.AddressInfo.Coordinates.Latitude,
					},
				},
			},
		}

		// Para shipments (pickup + delivery), incluir origen con información completa
		if hasPickup {
			order.Origin = domain.NodeInfo{
				ReferenceID: domain.ReferenceID(uuid.New().String()),
				AddressInfo: domain.AddressInfo{
					// Información de contacto del origen
					Contact: domain.Contact{
						FullName:     visit.Pickup.AddressInfo.Contact.FullName,
						PrimaryEmail: visit.Pickup.AddressInfo.Contact.Email,
						PrimaryPhone: visit.Pickup.AddressInfo.Contact.Phone,
						NationalID:   visit.Pickup.AddressInfo.Contact.NationalID,
					},
					// Información política/geográfica
					PoliticalArea: domain.PoliticalArea{
						Code:            visit.Pickup.AddressInfo.PoliticalArea.Code,
						AdminAreaLevel1: visit.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel1,
						AdminAreaLevel2: visit.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel2,
						AdminAreaLevel3: visit.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel3,
						AdminAreaLevel4: visit.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel4,
					},
					// Información de dirección
					AddressLine1: visit.Pickup.AddressInfo.AddressLine1,
					AddressLine2: visit.Pickup.AddressInfo.AddressLine2,
					ZipCode:      visit.Pickup.AddressInfo.ZipCode,
					// Coordenadas
					Coordinates: domain.Coordinates{
						Point: orb.Point{
							visit.Pickup.AddressInfo.Coordinates.Longitude,
							visit.Pickup.AddressInfo.Coordinates.Latitude,
						},
					},
				},
			}
		}

		// Mapear unidades de entrega con información completa
		deliveryUnits := make(domain.DeliveryUnits, 0, len(orderReq.DeliveryUnits))
		for _, du := range orderReq.DeliveryUnits {
			// Mapear los items de cada unidad de entrega
			items := make([]domain.Item, 0, len(du.Items))
			for _, item := range du.Items {
				items = append(items, domain.Item{
					Sku: item.Sku,
				})
			}

					// Crear la unidad de entrega del dominio con información completa
		// Create copies of the values to ensure we have valid pointers
		volume := du.Volume
		weight := du.Weight
		insurance := du.Insurance
		
		deliveryUnit := domain.DeliveryUnit{
			Lpn:       du.Lpn,
			Volume:    &volume,
			Weight:    &weight,
			Insurance: &insurance,
			Items:     items,
		}

			deliveryUnits = append(deliveryUnits, deliveryUnit)
		}
		order.DeliveryUnits = deliveryUnits

		// Mapear fechas de disponibilidad de recolección si está disponible
		if hasPickup && visit.Pickup.TimeWindow.Start != "" {
			order.CollectAvailabilityDate = domain.CollectAvailabilityDate{
				TimeRange: domain.TimeRange{
					StartTime: visit.Pickup.TimeWindow.Start,
					EndTime:   visit.Pickup.TimeWindow.End,
				},
			}
		}

		// Mapear fechas prometidas si está disponible
		if visit.Delivery.TimeWindow.Start != "" {
			order.PromisedDate = domain.PromisedDate{
				TimeRange: domain.TimeRange{
					StartTime: visit.Delivery.TimeWindow.Start,
					EndTime:   visit.Delivery.TimeWindow.End,
				},
			}
		}

		// Mapear instrucciones de entrega si está disponible
		if visit.Delivery.Instructions != "" {
			order.DeliveryInstructions = visit.Delivery.Instructions
		}

		// Mapear referencias si están disponibles
		if len(visit.Delivery.NodeInfo.ReferenceID) > 0 {
			order.References = []domain.Reference{
				{
					Type:  "node_reference",
					Value: visit.Delivery.NodeInfo.ReferenceID,
				},
			}
		}

		orders = append(orders, order)
	}

	return orders
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

// createOrderFromOptimizationOrder crea una orden del dominio desde una orden de optimización
func createOrderFromOptimizationOrder(optOrder optimization.Order) domain.Order {
	// Mapear las unidades de entrega con información completa
	deliveryUnits := make(domain.DeliveryUnits, 0, len(optOrder.DeliveryUnits))
	for _, du := range optOrder.DeliveryUnits {
		// Mapear los items de cada unidad de entrega
		items := make([]domain.Item, 0, len(du.Items))
		for _, item := range du.Items {
			items = append(items, domain.Item{
				Sku: item.Sku,
			})
		}

		// Crear la unidad de entrega del dominio con información completa
		// Create copies of the values to ensure we have valid pointers
		volume := du.Volume
		weight := du.Weight
		insurance := du.Insurance
		
		deliveryUnit := domain.DeliveryUnit{
			Lpn:       du.Lpn,
			Volume:    &volume,
			Weight:    &weight,
			Insurance: &insurance,
			Items:     items,
		}

		deliveryUnits = append(deliveryUnits, deliveryUnit)
	}

	// Crear la orden del dominio con información completa
	return domain.Order{
		ReferenceID:   domain.ReferenceID(optOrder.ReferenceID),
		DeliveryUnits: deliveryUnits,
	}
}
