package mapper

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/domain/optimization"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

// MapOptimizedFleetToPlan convierte un OptimizedFleet a un domain.Plan
func MapOptimizedFleetToPlan(
	ctx context.Context,
	optimizedFleet optimization.OptimizedFleet,
) domain.Plan {

	// Convertir rutas optimizadas a rutas del dominio
	domainRoutes := make([]domain.Route, 0, len(optimizedFleet.Routes))
	for _, optimizedRoute := range optimizedFleet.Routes {
		domainRoute := mapOptimizedRouteToDomainRoute(ctx, optimizedRoute)
		domainRoutes = append(domainRoutes, domainRoute)
	}

	// Convertir elementos no asignados
	unassignedOrders := make([]domain.Order, 0, len(optimizedFleet.Unassigned.Orders))
	for _, optOrder := range optimizedFleet.Unassigned.Orders {
		domainOrder := mapOptimizationOrderToDomainOrder(ctx, optOrder)
		unassignedOrders = append(unassignedOrders, domainOrder)
	}

	unassignedVehicles := make([]domain.Vehicle, 0, len(optimizedFleet.Unassigned.Vehicles))
	for _, optVehicle := range optimizedFleet.Unassigned.Vehicles {
		domainVehicle := mapOptimizationVehicleToDomainVehicle(ctx, optVehicle)
		unassignedVehicles = append(unassignedVehicles, domainVehicle)
	}

	unassignedOrigins := make([]domain.NodeInfo, 0, len(optimizedFleet.Unassigned.Origins))
	for _, optNodeInfo := range optimizedFleet.Unassigned.Origins {
		domainNodeInfo := mapOptimizationNodeInfoToDomainNodeInfo(ctx, optNodeInfo)
		unassignedOrigins = append(unassignedOrigins, domainNodeInfo)
	}

	return domain.Plan{
		ReferenceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		UnassignedOrigins:  unassignedOrigins,
		UnassignedVehicles: unassignedVehicles,
		UnassignedOrders:   unassignedOrders,
		PlannedDate:        optimizedFleet.PlannedDate,
		Routes:             domainRoutes,
	}
}

// mapOptimizedRouteToDomainRoute convierte una OptimizedRoute a domain.Route
func mapOptimizedRouteToDomainRoute(ctx context.Context, optRoute optimization.OptimizedRoute) domain.Route {
	// Para las rutas optimizadas, necesitamos extraer información de los pasos
	var origin, destination domain.NodeInfo
	var orders []domain.Order
	var vehicle domain.Vehicle

	// Buscar el primer paso con ubicación válida para el origen
	for _, step := range optRoute.Steps {
		if step.Location.Latitude != 0 || step.Location.Longitude != 0 {
			origin = mapOptimizationNodeInfoToDomainNodeInfo(ctx, step.Location.NodeInfo)
			break
		}
	}

	// Buscar el último paso con ubicación válida para el destino
	for i := len(optRoute.Steps) - 1; i >= 0; i-- {
		step := optRoute.Steps[i]
		if step.Location.Latitude != 0 || step.Location.Longitude != 0 {
			destination = mapOptimizationNodeInfoToDomainNodeInfo(ctx, step.Location.NodeInfo)
			break
		}
	}

	// Recolectar todas las órdenes de los pasos
	for _, step := range optRoute.Steps {
		for _, optOrder := range step.Orders {
			domainOrder := mapOptimizationOrderToDomainOrder(ctx, optOrder)
			orders = append(orders, domainOrder)
		}
	}

	// Crear un vehículo básico con la placa
	vehicle = domain.Vehicle{
		Plate: optRoute.VehiclePlate,
	}

	return domain.Route{
		ReferenceID: uuid.New().String(),
		Origin:      origin,
		Destination: destination,
		Vehicle:     vehicle,
		Orders:      orders,
	}
}

// mapOptimizationOrderToDomainOrder convierte una optimization.Order a domain.Order
func mapOptimizationOrderToDomainOrder(ctx context.Context, optOrder optimization.Order) domain.Order {
	deliveryUnits := make([]domain.DeliveryUnit, 0, len(optOrder.DeliveryUnits))
	for _, optDeliveryUnit := range optOrder.DeliveryUnits {
		domainDeliveryUnit := mapOptimizationDeliveryUnitToDomainDeliveryUnit(ctx, optDeliveryUnit)
		deliveryUnits = append(deliveryUnits, domainDeliveryUnit)
	}

	return domain.Order{
		ReferenceID:   domain.ReferenceID(optOrder.ReferenceID),
		DeliveryUnits: deliveryUnits,
		// Nota: Algunos campos del domain.Order no están disponibles en optimization.Order
		// Se pueden establecer valores por defecto o requerir parámetros adicionales
	}
}

// mapOptimizationDeliveryUnitToDomainDeliveryUnit convierte una optimization.DeliveryUnit a domain.DeliveryUnit
func mapOptimizationDeliveryUnitToDomainDeliveryUnit(ctx context.Context, optDeliveryUnit optimization.DeliveryUnit) domain.DeliveryUnit {
	items := make([]domain.Item, 0, len(optDeliveryUnit.Items))
	for _, optItem := range optDeliveryUnit.Items {
		domainItem := mapOptimizationItemToDomainItem(ctx, optItem)
		items = append(items, domainItem)
	}

	return domain.DeliveryUnit{
		Lpn:       optDeliveryUnit.Lpn,
		Weight:    domain.Weight{Value: optDeliveryUnit.Weight},
		Insurance: domain.Insurance{UnitValue: optDeliveryUnit.Insurance},
		Items:     items,
		// Nota: Algunos campos del domain.DeliveryUnit no están disponibles en optimization.DeliveryUnit
	}
}

// mapOptimizationItemToDomainItem convierte una optimization.Item a domain.Item
func mapOptimizationItemToDomainItem(ctx context.Context, optItem optimization.Item) domain.Item {
	return domain.Item{
		Sku: optItem.Sku,
		// Nota: Algunos campos del domain.Item no están disponibles en optimization.Item
	}
}

// mapOptimizationVehicleToDomainVehicle convierte una optimization.Vehicle a domain.Vehicle
func mapOptimizationVehicleToDomainVehicle(ctx context.Context, optVehicle optimization.Vehicle) domain.Vehicle {
	return domain.Vehicle{
		Plate: optVehicle.Plate,
		// Nota: Algunos campos del domain.Vehicle no están disponibles en optimization.Vehicle
	}
}

// mapOptimizationNodeInfoToDomainNodeInfo convierte una optimization.NodeInfo a domain.NodeInfo
func mapOptimizationNodeInfoToDomainNodeInfo(ctx context.Context, optNodeInfo optimization.NodeInfo) domain.NodeInfo {
	return domain.NodeInfo{
		ReferenceID: domain.ReferenceID(optNodeInfo.ReferenceID),
		// Nota: Algunos campos del domain.NodeInfo no están disponibles en optimization.NodeInfo
	}
}
