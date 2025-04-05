package tidbrepository

import (
	"context"
	"errors"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertPlan func(context.Context, domain.Plan) error

func init() {
	ioc.Registry(
		NewUpsertPlan,
		tidb.NewTIDBConnection,
		NewLoadOrderStatuses)
}

func NewUpsertPlan(conn tidb.TIDBConnection, loadOrderStatuses LoadOrderStatuses) UpsertPlan {
	return func(ctx context.Context, p domain.Plan) error {
		plannedStatusID := loadOrderStatuses().Planned().ID
		availableStatusID := loadOrderStatuses().Available().ID
		err := conn.Transaction(func(tx *gorm.DB) error {
			// 1️⃣ Buscar o crear el plan
			var plan table.Plan
			err := tx.Table("plans").
				Where("reference_id = ? AND organization_id = ?",
					p.ReferenceID, p.Organization.ID).
				First(&plan).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				plan = mapper.MapPlan(p)
				plan.CreatedAt = time.Now()
				plan.OrganizationID = p.Organization.ID
				err = tx.Table("plans").Create(&plan).Error
			} else if err == nil {
				planToUpdate := plan.Map().UpdateIfChanged(p)
				planToCreate := mapper.MapPlan(planToUpdate)
				planToCreate.CreatedAt = plan.CreatedAt
				err = tx.Table("plans").Save(&planToCreate).Error
				plan = planToCreate
			}
			if err != nil {
				return err
			}

			err = tx.Table("orders").
				Where("plan_id = ?", plan.ID).
				Updates(map[string]interface{}{
					"route_id":        nil,
					"plan_id":         nil,
					"order_status_id": availableStatusID,
				}).Error
			if err != nil {
				return err
			}

			// 2️⃣ Obtener `reference_id` de todas las rutas del plan
			var routeReferenceIDs []string
			routeSet := make(map[string]bool) // Evita duplicados
			for _, route := range p.Routes {
				if !routeSet[route.ReferenceID] {
					routeReferenceIDs = append(routeReferenceIDs, route.ReferenceID)
					routeSet[route.ReferenceID] = true
				}
			}

			// 3️⃣ Consultar rutas existentes
			var existingRoutes []table.Route
			err = tx.Table("routes").
				Where("reference_id IN ? AND organization_id = ?",
					routeReferenceIDs, p.Organization.ID).
				Find(&existingRoutes).Error
			if err != nil {
				return err
			}

			// 4️⃣ Crear un mapa de rutas existentes
			existingRouteMap := make(map[string]int64)
			for _, route := range existingRoutes {
				existingRouteMap[route.ReferenceID] = route.ID
			}

			// 5️⃣ Insertar rutas que no existan
			var newRoutes []table.Route
			for _, route := range p.Routes {
				if _, exists := existingRouteMap[route.ReferenceID]; !exists {
					newRoute := mapper.MapRoute(route)
					newRoute.PlanID = plan.ID
					newRoute.OrganizationID = p.Organization.ID
					newRoutes = append(newRoutes, newRoute)
				}
			}
			if len(newRoutes) > 0 {
				err = tx.Table("routes").Create(&newRoutes).Error
				if err != nil {
					return err
				}
				// Actualizar el mapa de rutas con los IDs recién insertados
				for _, route := range newRoutes {
					existingRouteMap[route.ReferenceID] = route.ID
				}
			}

			// 6️⃣ Consolidar todas las órdenes
			var allOrders []domain.Order
			var referenceIDs []string
			routeIDMap := make(map[string]*int64)

			orderSet := make(map[string]bool) // Evita duplicados
			for _, route := range p.Routes {
				for _, order := range route.Orders {
					if !orderSet[string(order.ReferenceID)] {
						allOrders = append(allOrders, order)
						referenceIDs = append(referenceIDs, string(order.ReferenceID))
						orderSet[string(order.ReferenceID)] = true
					}
					if routeID, exists := existingRouteMap[route.ReferenceID]; exists {
						routeIDMap[string(order.ReferenceID)] = &routeID
					}
				}
			}
			for _, order := range p.UnassignedOrders {
				if !orderSet[string(order.ReferenceID)] {
					allOrders = append(allOrders, order)
					referenceIDs = append(referenceIDs, string(order.ReferenceID))
					orderSet[string(order.ReferenceID)] = true
				}
			}

			// 7️⃣ Obtener todas las órdenes existentes
			var existingOrders []table.Order
			err = tx.Table("orders").
				Where("reference_id IN ? AND organization_id = ?",
					referenceIDs, p.Organization.ID).
				Find(&existingOrders).Error
			if err != nil {
				return err
			}

			// 8️⃣ Crear un mapa de órdenes existentes
			existingOrdersMap := make(map[string]bool)
			for _, existingOrder := range existingOrders {
				existingOrdersMap[existingOrder.ReferenceID] = true
			}

			// 9️⃣ Insertar órdenes en batch solo si no existen
			var newOrders []table.Order
			for _, order := range allOrders {
				if existingOrdersMap[string(order.ReferenceID)] {
					continue
				}
				newOrder := mapper.MapOrderToTable(order)
				//newOrder.PlanID = &plan.ID
				//newOrder.OrderStatusID = plannedStatusID
				newOrder.OrganizationID = p.Organization.ID
				/*
					if routeID, exists := routeIDMap[string(order.ReferenceID)]; exists {
						newOrder.RouteID = routeID
					}
				*/
				newOrders = append(newOrders, newOrder)
			}
			if len(newOrders) > 0 {
				err = tx.Table("orders").Create(&newOrders).Error
				if err != nil {
					return err
				}
			}

			// 🔟 Actualizar en batch la data del plan para todas las órdenes
			for _, order := range allOrders {
				plannedData := table.JSONPlannedData{
					JSONPlanLocation: table.PlanLocation{
						Longitude: order.Destination.AddressInfo.Location.Lat(),
						Latitude:  order.Destination.AddressInfo.Location.Lon(),
					},
					JSONPlanCorrectedLocation: table.PlanLocation{
						Longitude: order.Destination.AddressInfo.CorrectedLocation.Lat(),
						Latitude:  order.Destination.AddressInfo.CorrectedLocation.Lon(),
					},
					PlanCorrectedDistance: order.Destination.AddressInfo.CorrectedDistance,
				}

				updateFields := map[string]interface{}{
					"plan_id":           plan.ID,
					"order_status_id":   plannedStatusID,
					"json_planned_data": plannedData,
					"sequence_number":   order.SequenceNumber,
					"organization_id":   p.Organization.ID,
				}

				if routeID, exists := routeIDMap[string(order.ReferenceID)]; exists {
					updateFields["route_id"] = *routeID
				} else {
					updateFields["route_id"] = nil
				}

				err = tx.Table("orders").
					Where("reference_id = ? AND organization_id = ?",
						order.ReferenceID, p.Organization.ID).
					Updates(updateFields).Error
				if err != nil {
					return err
				}
			}

			return nil
		})

		return err
	}
}
