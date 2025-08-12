package vroom

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
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

type Optimize func(ctx context.Context, request optimization.FleetOptimization) ([]request.UpsertRouteRequest, error)

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
	return func(ctx context.Context, fleetOptimization optimization.FleetOptimization) ([]request.UpsertRouteRequest, error) {

		vroomRequest, err := mapper.MapOptimizationRequest(ctx, fleetOptimization)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		if res.IsError() {
			obs.Logger.ErrorContext(ctx,
				"VROOM_API_ERROR",
				"status", res.StatusCode(),
				"body", res.String(),
				"request", vroomRequest,
			)

			return nil, fmt.Errorf("VROOM API error (status %d): %s\nRequest payload: %+v",
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
			return nil, fmt.Errorf("failed to deserialize VROOM response: %w", err)
		}

		planReferenceID := uuid.New().String()

		// Slice para almacenar todos los polylines consolidados
		var allPolylines []string
		var allRouteData []model.RouteData
		var routeRequests []request.UpsertRouteRequest

		fleetOptimizations, unassignedOrders := vroomResponse.MapOptimizationRequests(ctx, fleetOptimization)

		// Crear un mapa de motivos de no asignación
		unassignedReasons := make(map[string]string)

		// Extraer motivos de la respuesta de VROOM
		obs.Logger.InfoContext(ctx, "VROOM_UNASSIGNED_JOBS", "count", len(vroomResponse.Unassigned))
		for _, unassigned := range vroomResponse.Unassigned {
			obs.Logger.InfoContext(ctx, "VROOM_UNASSIGNED_JOB",
				"id", unassigned.ID,
				"reason", unassigned.Reason,
				"location", unassigned.Location)

			// Buscar la visita correspondiente al job sin asignar
			if int(unassigned.ID) <= len(fleetOptimization.Visits) {
				visit := fleetOptimization.Visits[unassigned.ID-1]
				// Asociar el motivo con todas las órdenes de esta visita
				for _, order := range visit.Orders {
					unassignedReasons[order.ReferenceID] = unassigned.Reason
					obs.Logger.InfoContext(ctx, "MAPPED_UNASSIGNED_REASON",
						"orderID", order.ReferenceID,
						"reason", unassigned.Reason)
				}
			}
		}

		// Crear ruta especial para órdenes sin asignar si las hay
		if len(unassignedOrders) > 0 {
			// Para órdenes sin asignar, crear un vehículo vacío
			emptyVehicle := optimization.Vehicle{
				Plate: "UNASSIGNED",
				StartLocation: optimization.VehicleLocation{
					AddressInfo: optimization.AddressInfo{},
					NodeInfo:    optimization.NodeInfo{},
				},
				EndLocation: optimization.VehicleLocation{
					AddressInfo: optimization.AddressInfo{},
					NodeInfo:    optimization.NodeInfo{},
				},
				Skills:     []string{},
				TimeWindow: optimization.TimeWindow{},
				Capacity:   optimization.Capacity{},
			}
			// Para órdenes sin asignar, usar las visitas originales
			unassignedRouteRequest := createUnassignedRouteRequest(unassignedOrders, planReferenceID, emptyVehicle, fleetOptimization.Visits, unassignedReasons)
			routeRequests = append(routeRequests, unassignedRouteRequest)
		}

		// Segunda optimización - optimización individual por ruta
		for optimizationIndex, individualFleetOptimization := range fleetOptimizations {
			individualVroomRequest, err := mapper.MapOptimizationRequest(ctx, individualFleetOptimization)
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Failed to map individual optimization request", "error", err)
				continue
			}

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

				// Convertir la ruta del dominio a UpsertRouteRequest
				// Usar la información del vehículo original del input
				var originalVehicle optimization.Vehicle
				if len(individualFleetOptimization.Vehicles) > 0 {
					originalVehicle = individualFleetOptimization.Vehicles[0]
				}
				routeRequest := createUpsertRouteRequest(route, planReferenceID, originalVehicle, individualFleetOptimization.Visits)
				routeRequests = append(routeRequests, routeRequest)

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
		}

		return routeRequests, nil
	}
}

// createUpsertRouteRequest convierte una ruta del dominio a UpsertRouteRequest
func createUpsertRouteRequest(route domain.Route, planReferenceID string, originalVehicle optimization.Vehicle, originalVisits []optimization.Visit) request.UpsertRouteRequest {
	// Agrupar órdenes por secuencia, dirección y contacto
	visitGroups := groupOrdersByVisit(route.Orders)

	// Convertir grupos a visitas usando información de las visitas originales
	visits := make([]request.UpsertRouteVisit, 0, len(visitGroups))
	for _, group := range visitGroups {
		// Buscar la visita original correspondiente
		var originalVisit *optimization.Visit
		for _, origVisit := range originalVisits {
			// Buscar por dirección de delivery
			if origVisit.Delivery.AddressInfo.Coordinates.Latitude == group.AddressInfo.Coordinates.Point[1] &&
				origVisit.Delivery.AddressInfo.Coordinates.Longitude == group.AddressInfo.Coordinates.Point[0] {
				originalVisit = &origVisit
				break
			}
		}

		visit := mapOrderGroupToVisitFromOptimization(group, originalVisit)
		visits = append(visits, visit)
	}

	// Ordenar visitas por número de secuencia
	sort.Slice(visits, func(i, j int) bool {
		return visits[i].SequenceNumber < visits[j].SequenceNumber
	})

	return request.UpsertRouteRequest{
		ReferenceID:     route.ReferenceID,
		PlanReferenceID: planReferenceID,
		Vehicle:         mapVehicleToRequestFromOptimization(originalVehicle),
		Geometry: request.UpsertRouteGeometry{
			Encoding: route.Geometry.Encoding,
			Type:     route.Geometry.Type,
			Value:    route.Geometry.Value,
		},
		Visits:    visits,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// OrderGroup representa un grupo de órdenes que se pueden agrupar en una visita
type OrderGroup struct {
	SequenceNumber       int
	AddressInfo          domain.AddressInfo
	Contact              domain.Contact
	DeliveryInstructions string
	Orders               []domain.Order
}

// groupOrdersByVisit agrupa las órdenes por secuencia, dirección y contacto
func groupOrdersByVisit(orders []domain.Order) []OrderGroup {
	if len(orders) == 0 {
		return []OrderGroup{}
	}

	// Crear un mapa para agrupar órdenes
	groups := make(map[string]OrderGroup)
	sequenceCounter := 1

	for _, order := range orders {
		// Crear una clave única basada en la dirección y contacto del destino
		destKey := createDestinationKey(order.Destination.AddressInfo)

		if group, exists := groups[destKey]; exists {
			// Agregar la orden al grupo existente
			group.Orders = append(group.Orders, order)
			groups[destKey] = group
		} else {
			// Crear un nuevo grupo
			groups[destKey] = OrderGroup{
				SequenceNumber:       sequenceCounter,
				AddressInfo:          order.Destination.AddressInfo,
				Contact:              order.Destination.AddressInfo.Contact,
				DeliveryInstructions: order.DeliveryInstructions,
				Orders:               []domain.Order{order},
			}
			sequenceCounter++
		}
	}

	// Convertir el mapa a slice
	result := make([]OrderGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// createDestinationKey crea una clave única para agrupar órdenes por destino
func createDestinationKey(addrInfo domain.AddressInfo) string {
	// Usar coordenadas como clave principal
	lat := addrInfo.Coordinates.Point[1]
	lon := addrInfo.Coordinates.Point[0]

	// También incluir información de contacto para mayor precisión
	contactKey := fmt.Sprintf("%s-%s", addrInfo.Contact.FullName, addrInfo.Contact.PrimaryEmail)

	return fmt.Sprintf("%.6f-%.6f-%s", lat, lon, contactKey)
}

// mapOrderGroupToVisit convierte un grupo de órdenes a una visita
func mapOrderGroupToVisit(group OrderGroup) request.UpsertRouteVisit {
	// Mapear órdenes del grupo
	orders := make([]request.UpsertRouteOrder, 0, len(group.Orders))
	for _, order := range group.Orders {
		modelOrder := mapOrderToRequest(order)
		orders = append(orders, modelOrder)
	}

	// Buscar información de TimeWindow y ServiceTime de la primera orden
	var timeWindow request.UpsertRouteTimeWindow
	var serviceTime int64 = 0

	if len(group.Orders) > 0 {
		// Intentar obtener TimeWindow de la primera orden
		// Por ahora usamos valores por defecto, se puede implementar mapeo específico
		timeWindow = request.UpsertRouteTimeWindow{
			Start: "", // Por defecto vacío, se puede implementar mapeo específico
			End:   "",
		}
	}

	return request.UpsertRouteVisit{
		Type:                 "delivery",
		Instructions:         "", // Por defecto vacío, se puede implementar mapeo específico
		AddressInfo:          mapAddressInfoToRequest(group.AddressInfo),
		NodeInfo:             mapNodeInfoToRequest(group.AddressInfo),
		DeliveryInstructions: group.DeliveryInstructions,
		SequenceNumber:       group.SequenceNumber,
		ServiceTime:          serviceTime,
		TimeWindow:           timeWindow,
		Orders:               orders,
	}
}

// mapOrderGroupToVisitFromOptimization convierte un grupo de órdenes de optimización a una visita
func mapOrderGroupToVisitFromOptimization(group OrderGroup, originalVisit *optimization.Visit) request.UpsertRouteVisit {
	// Mapear órdenes del grupo usando la información de optimización
	orders := make([]request.UpsertRouteOrder, 0, len(group.Orders))
	for _, order := range group.Orders {
		// Buscar la orden correspondiente en la visita original
		var originalOrder optimization.Order
		if originalVisit != nil {
			for _, origOrder := range originalVisit.Orders {
				if origOrder.ReferenceID == order.ReferenceID.String() {
					originalOrder = origOrder
					break
				}
			}
		}

		if originalOrder.ReferenceID != "" {
			modelOrder := mapOrderToRequestFromOptimization(originalOrder)
			orders = append(orders, modelOrder)
		} else {
			// Fallback a mapeo del dominio si no se encuentra
			modelOrder := mapOrderToRequest(order)
			orders = append(orders, modelOrder)
		}
	}

	// Preparar valores por defecto
	instructions := ""
	serviceTime := int64(0)
	timeWindow := request.UpsertRouteTimeWindow{
		Start: "",
		End:   "",
	}
	nodeInfo := request.UpsertRouteNodeInfo{
		ReferenceID: "",
	}

	// Usar información de la visita original si está disponible
	if originalVisit != nil {
		instructions = originalVisit.Delivery.Instructions
		serviceTime = originalVisit.Delivery.ServiceTime
		timeWindow = request.UpsertRouteTimeWindow{
			Start: originalVisit.Delivery.TimeWindow.Start,
			End:   originalVisit.Delivery.TimeWindow.End,
		}
		nodeInfo = request.UpsertRouteNodeInfo{
			ReferenceID: originalVisit.Delivery.NodeInfo.ReferenceID,
		}
	}

	return request.UpsertRouteVisit{
		Type:                 "delivery",
		Instructions:         instructions,
		AddressInfo:          mapAddressInfoToRequest(group.AddressInfo),
		NodeInfo:             nodeInfo,
		DeliveryInstructions: group.DeliveryInstructions,
		SequenceNumber:       group.SequenceNumber,
		ServiceTime:          serviceTime,
		TimeWindow:           timeWindow,
		Orders:               orders,
	}
}

// mapOrderToRequest convierte una orden del dominio al request
func mapOrderToRequest(order domain.Order) request.UpsertRouteOrder {
	deliveryUnits := make([]request.UpsertRouteDeliveryUnit, 0, len(order.DeliveryUnits))
	for _, du := range order.DeliveryUnits {
		modelDU := mapDeliveryUnitToRequest(du)
		deliveryUnits = append(deliveryUnits, modelDU)
	}

	return request.UpsertRouteOrder{
		ReferenceID:   order.ReferenceID.String(),
		DeliveryUnits: deliveryUnits,
	}
}

// mapOrderToRequestFromOptimization convierte una orden de optimización al request
func mapOrderToRequestFromOptimization(order optimization.Order) request.UpsertRouteOrder {
	deliveryUnits := make([]request.UpsertRouteDeliveryUnit, 0, len(order.DeliveryUnits))
	for _, du := range order.DeliveryUnits {
		modelDU := mapDeliveryUnitToRequestFromOptimization(du)
		deliveryUnits = append(deliveryUnits, modelDU)
	}

	return request.UpsertRouteOrder{
		ReferenceID:   order.ReferenceID,
		DeliveryUnits: deliveryUnits,
	}
}

// mapDeliveryUnitToRequestFromOptimization convierte una unidad de entrega de optimización al request
func mapDeliveryUnitToRequestFromOptimization(du optimization.DeliveryUnit) request.UpsertRouteDeliveryUnit {
	// Mapear items
	items := make([]request.UpsertRouteItem, 0, len(du.Items))
	for _, item := range du.Items {
		items = append(items, request.UpsertRouteItem{
			Sku:         item.Sku,
			Description: item.Description,
			Quantity:    item.Quantity,
		})
	}

	return request.UpsertRouteDeliveryUnit{
		Items:  items,
		Volume: du.Volume,
		Weight: du.Weight,
		Price:  du.Price,
		Lpn:    du.Lpn,
		Skills: du.Skills, // Incluir skills desde la optimización
	}
}

// mapDeliveryUnitToRequest convierte una unidad de entrega del dominio al request
func mapDeliveryUnitToRequest(du domain.DeliveryUnit) request.UpsertRouteDeliveryUnit {
	// Mapear items
	items := make([]request.UpsertRouteItem, 0, len(du.Items))
	for _, item := range du.Items {
		items = append(items, request.UpsertRouteItem{
			Sku: item.Sku,
		})
	}

	// Mapear skills (por ahora vacío, se puede implementar mapeo específico desde el dominio)
	skills := make([]string, 0)

	return request.UpsertRouteDeliveryUnit{
		Items:  items,
		Volume: int64(*du.Volume),
		Weight: int64(*du.Weight),
		Price:  int64(*du.Price),
		Lpn:    du.Lpn,
		Skills: skills,
	}
}

// mapVehicleToRequestFromOptimization convierte el vehículo de optimización al request
func mapVehicleToRequestFromOptimization(vehicle optimization.Vehicle) request.UpsertRouteVehicle {
	// Mapear StartLocation con información completa
	startLocation := request.UpsertRouteVehicleLocation{
		AddressInfo: mapAddressInfoToRequestFromOptimization(vehicle.StartLocation.AddressInfo),
		NodeInfo: request.UpsertRouteNodeInfo{
			ReferenceID: vehicle.StartLocation.NodeInfo.ReferenceID,
		},
	}

	// Mapear EndLocation con información completa
	endLocation := request.UpsertRouteVehicleLocation{
		AddressInfo: mapAddressInfoToRequestFromOptimization(vehicle.EndLocation.AddressInfo),
		NodeInfo: request.UpsertRouteNodeInfo{
			ReferenceID: vehicle.EndLocation.NodeInfo.ReferenceID,
		},
	}

	// Mapear TimeWindow
	timeWindow := request.UpsertRouteTimeWindow{
		Start: vehicle.TimeWindow.Start,
		End:   vehicle.TimeWindow.End,
	}

	return request.UpsertRouteVehicle{
		Plate:         vehicle.Plate,
		StartLocation: startLocation,
		EndLocation:   endLocation,
		Skills:        vehicle.Skills,
		TimeWindow:    timeWindow,
		Capacity: request.UpsertRouteVehicleCapacity{
			Volume:                vehicle.Capacity.Volume,
			Weight:                vehicle.Capacity.Weight,
			Insurance:             vehicle.Capacity.Insurance,
			DeliveryUnitsQuantity: vehicle.Capacity.DeliveryUnitsQuantity,
		},
	}
}

// mapAddressInfoToRequestFromOptimization convierte AddressInfo de optimización al request
func mapAddressInfoToRequestFromOptimization(addrInfo optimization.AddressInfo) request.UpsertRouteAddressInfo {
	return request.UpsertRouteAddressInfo{
		AddressLine1:  addrInfo.AddressLine1,
		AddressLine2:  addrInfo.AddressLine2,
		Contact:       mapContactToRequestFromOptimization(addrInfo.Contact),
		Coordinates:   mapCoordinatesToRequestFromOptimization(addrInfo.Coordinates),
		PoliticalArea: mapPoliticalAreaToRequestFromOptimization(addrInfo.PoliticalArea),
		ZipCode:       addrInfo.ZipCode,
	}
}

// mapContactToRequestFromOptimization convierte Contact de optimización al request
func mapContactToRequestFromOptimization(contact optimization.Contact) request.UpsertRouteContact {
	return request.UpsertRouteContact{
		Email:      contact.Email,
		Phone:      contact.Phone,
		NationalID: contact.NationalID,
		FullName:   contact.FullName,
	}
}

// mapCoordinatesToRequestFromOptimization convierte Coordinates de optimización al request
func mapCoordinatesToRequestFromOptimization(coords optimization.Coordinates) request.UpsertRouteCoordinates {
	return request.UpsertRouteCoordinates{
		Latitude:  coords.Latitude,
		Longitude: coords.Longitude,
	}
}

// mapPoliticalAreaToRequestFromOptimization convierte PoliticalArea de optimización al request
func mapPoliticalAreaToRequestFromOptimization(pa optimization.PoliticalArea) request.UpsertRoutePoliticalArea {
	return request.UpsertRoutePoliticalArea{
		Code:            pa.Code,
		AdminAreaLevel1: pa.AdminAreaLevel1,
		AdminAreaLevel2: pa.AdminAreaLevel2,
		AdminAreaLevel3: pa.AdminAreaLevel3,
		AdminAreaLevel4: pa.AdminAreaLevel4,
	}
}

// mapVehicleToRequest convierte el vehículo del dominio al request
func mapVehicleToRequest(vehicle domain.Vehicle) request.UpsertRouteVehicle {
	// Mapear StartLocation si está disponible
	startLocation := request.UpsertRouteVehicleLocation{
		AddressInfo: request.UpsertRouteAddressInfo{}, // Por defecto vacío
		NodeInfo:    request.UpsertRouteNodeInfo{},    // Por defecto vacío
	}

	// Mapear EndLocation si está disponible
	endLocation := request.UpsertRouteVehicleLocation{
		AddressInfo: request.UpsertRouteAddressInfo{}, // Por defecto vacío
		NodeInfo:    request.UpsertRouteNodeInfo{},    // Por defecto vacío
	}

	// Mapear TimeWindow si está disponible
	timeWindow := request.UpsertRouteTimeWindow{
		Start: "", // Por defecto vacío
		End:   "",
	}

	return request.UpsertRouteVehicle{
		Plate:         vehicle.Plate,
		StartLocation: startLocation,
		EndLocation:   endLocation,
		Skills:        []string{}, // Por defecto vacío, se puede implementar mapeo específico
		TimeWindow:    timeWindow,
		Capacity: request.UpsertRouteVehicleCapacity{
			Volume:                int64(vehicle.Weight.Value), // Usar Weight.Value como volumen
			Weight:                int64(vehicle.Weight.Value),
			Insurance:             int64(vehicle.Insurance.MaxInsuranceCoverage.Amount),
			DeliveryUnitsQuantity: 0, // Por defecto 0, se puede implementar mapeo específico
		},
	}
}

// mapAddressInfoToRequest convierte AddressInfo del dominio al request
func mapAddressInfoToRequest(addr domain.AddressInfo) request.UpsertRouteAddressInfo {
	return request.UpsertRouteAddressInfo{
		AddressLine1:  addr.AddressLine1,
		AddressLine2:  addr.AddressLine2,
		Contact:       mapContactToRequest(addr.Contact),
		Coordinates:   mapCoordinatesToRequest(addr.Coordinates),
		PoliticalArea: mapPoliticalAreaToRequest(addr.PoliticalArea),
		ZipCode:       addr.ZipCode,
	}
}

// mapContactToRequest convierte Contact del dominio al request
func mapContactToRequest(contact domain.Contact) request.UpsertRouteContact {
	return request.UpsertRouteContact{
		Email:      contact.PrimaryEmail,
		FullName:   contact.FullName,
		NationalID: contact.NationalID,
		Phone:      contact.PrimaryPhone,
	}
}

// mapCoordinatesToRequest convierte Coordinates del dominio al request
func mapCoordinatesToRequest(coords domain.Coordinates) request.UpsertRouteCoordinates {
	return request.UpsertRouteCoordinates{
		Latitude:  coords.Point[1], // orb.Point es [longitude, latitude]
		Longitude: coords.Point[0],
	}
}

// mapPoliticalAreaToRequest convierte PoliticalArea del dominio al request
func mapPoliticalAreaToRequest(pa domain.PoliticalArea) request.UpsertRoutePoliticalArea {
	return request.UpsertRoutePoliticalArea{
		Code:            pa.Code,
		AdminAreaLevel1: pa.AdminAreaLevel1,
		AdminAreaLevel2: pa.AdminAreaLevel2,
		AdminAreaLevel3: pa.AdminAreaLevel3,
		AdminAreaLevel4: pa.AdminAreaLevel4,
	}
}

// mapNodeInfoToRequest convierte la información del nodo del dominio al request
func mapNodeInfoToRequest(addr domain.AddressInfo) request.UpsertRouteNodeInfo {
	// Por ahora usar un ID vacío, se puede implementar lógica específica si es necesario
	return request.UpsertRouteNodeInfo{
		ReferenceID: "", // TODO: Implementar si es necesario
	}
}

// createUnassignedRouteRequest crea una ruta especial para órdenes sin asignar
func createUnassignedRouteRequest(unassignedOrders []optimization.Order, planReferenceID string, vehicle optimization.Vehicle, originalVisits []optimization.Visit, unassignedReasons map[string]string) request.UpsertRouteRequest {
	// Crear un vehículo "UNASSIGNED" con información mínima
	unassignedVehicle := request.UpsertRouteVehicle{
		Plate: "UNASSIGNED",
		StartLocation: request.UpsertRouteVehicleLocation{
			AddressInfo: request.UpsertRouteAddressInfo{},
			NodeInfo:    request.UpsertRouteNodeInfo{},
		},
		EndLocation: request.UpsertRouteVehicleLocation{
			AddressInfo: request.UpsertRouteAddressInfo{},
			NodeInfo:    request.UpsertRouteNodeInfo{},
		},
		Skills:     nil,
		TimeWindow: request.UpsertRouteTimeWindow{},
		Capacity: request.UpsertRouteVehicleCapacity{
			Volume:                0,
			Weight:                0,
			Insurance:             0,
			DeliveryUnitsQuantity: 0,
		},
	}

	// Crear visitas para las órdenes sin asignar
	var visits []request.UpsertRouteVisit

	for _, order := range unassignedOrders {
		// Buscar la visita original que contiene esta orden
		var originalVisit *optimization.Visit
		for _, visit := range originalVisits {
			for _, visitOrder := range visit.Orders {
				if visitOrder.ReferenceID == order.ReferenceID {
					originalVisit = &visit
					break
				}
			}
			if originalVisit != nil {
				break
			}
		}

		if originalVisit != nil {
			// Obtener el motivo de no asignación si está disponible
			unassignedReason := ""
			if reason, exists := unassignedReasons[order.ReferenceID]; exists {
				unassignedReason = reason
			}

			// Log para debug
			fmt.Printf("DEBUG: Orden %s - Motivo: '%s'\n", order.ReferenceID, unassignedReason)

			// Usar la información de la visita original
			visit := request.UpsertRouteVisit{
				Type:         "delivery",
				Instructions: originalVisit.Delivery.Instructions,
				AddressInfo:  mapAddressInfoToRequestFromOptimization(originalVisit.Delivery.AddressInfo),
				NodeInfo: request.UpsertRouteNodeInfo{
					ReferenceID: originalVisit.Delivery.NodeInfo.ReferenceID,
				},
				DeliveryInstructions: originalVisit.Delivery.Instructions,
				SequenceNumber:       1,
				ServiceTime:          originalVisit.Delivery.ServiceTime,
				TimeWindow: request.UpsertRouteTimeWindow{
					Start: originalVisit.Delivery.TimeWindow.Start,
					End:   originalVisit.Delivery.TimeWindow.End,
				},
				Orders: []request.UpsertRouteOrder{
					{
						ReferenceID: order.ReferenceID,
						DeliveryUnits: func() []request.UpsertRouteDeliveryUnit {
							units := make([]request.UpsertRouteDeliveryUnit, 0, len(order.DeliveryUnits))
							for _, du := range order.DeliveryUnits {
								units = append(units, mapDeliveryUnitToRequestFromOptimization(du))
							}
							return units
						}(),
					},
				},
				UnassignedReason: unassignedReason,
			}

			visits = append(visits, visit)
		}
	}

	return request.UpsertRouteRequest{
		ReferenceID:     fmt.Sprintf("UNASSIGNED-%s", uuid.New().String()),
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
		PlanReferenceID: planReferenceID,
		Vehicle:         unassignedVehicle,
		Geometry: request.UpsertRouteGeometry{
			Encoding: "",
			Type:     "",
			Value:    "",
		},
		Visits: visits,
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
			insurance := du.Price

			deliveryUnit := domain.DeliveryUnit{
				Lpn:    du.Lpn,
				Volume: &volume,
				Weight: &weight,
				Price:  &insurance,
				Items:  items,
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
		insurance := du.Price

		deliveryUnit := domain.DeliveryUnit{
			Lpn:    du.Lpn,
			Volume: &volume,
			Weight: &weight,
			Price:  &insurance,
			Items:  items,
			// Skills se mapean desde la optimización
		}

		deliveryUnits = append(deliveryUnits, deliveryUnit)
	}

	// Crear la orden del dominio con información completa
	return domain.Order{
		ReferenceID:   domain.ReferenceID(optOrder.ReferenceID),
		DeliveryUnits: deliveryUnits,
	}
}
