package tidbrepository

import (
	"context"
	"errors"
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

func NewUpsertPlan(conn tidb.TIDBConnection, loadOrdorderStatuses LoadOrderStatuses) UpsertPlan {
	return func(ctx context.Context, p domain.Plan) error {
		plannedStatusID := loadOrdorderStatuses().Planned().ID
		plan := table.Plan{}
		err := conn.DB.WithContext(ctx).Table("plans").
			Where("reference_id = ? AND organization_country_id = ?",
				p.ReferenceID, p.Organization.OrganizationCountryID).First(&plan).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		planToUpdate := plan.Map().UpdateIfChanged(p)
		DBPlanToUpdate := mapper.MapPlan(planToUpdate)
		DBPlanToUpdate.CreatedAt = plan.CreatedAt

		var routes []table.Route
		err = conn.DB.WithContext(ctx).Table("routes").
			Where("reference_id = ? AND organization_country_id = ?",
				p.ReferenceID, p.Organization.OrganizationCountryID).Find(&routes).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Mapeamos las rutas existentes por su referencia
		routeMap := make(map[string]table.Route)
		for _, r := range routes {
			routeMap[r.ReferenceID] = r
		}

		var DBRoutesToUpdate []table.Route

		// Aquí agrupamos las órdenes según la ruta, utilizando la estructura routeOrders
		type routeOrders struct {
			RouteReferenceID string
			Orders           []domain.Order
		}
		var routeOrdersMap []routeOrders
		var unassignedOrders []domain.Order

		// Procesamos cada ruta del plan
		for _, inputRoute := range p.Routes {
			var routeToUpdate domain.Route
			existingRoute, exists := routeMap[inputRoute.ReferenceID]
			if exists {
				routeToUpdate = existingRoute.Map().UpdateIfChanged(inputRoute)
			} else {
				routeToUpdate = inputRoute
			}

			DBRouteToUpdate := mapper.MapRoute(routeToUpdate)
			DBRouteToUpdate.PlanID = DBPlanToUpdate.ID
			if exists {
				DBRouteToUpdate.CreatedAt = existingRoute.CreatedAt
			}
			DBRoutesToUpdate = append(DBRoutesToUpdate, DBRouteToUpdate)

			// Agrupamos las órdenes de esta ruta
			var ordersForRoute []domain.Order
			for _, order := range inputRoute.Orders {
				ordersForRoute = append(ordersForRoute, order)
			}
			routeOrdersMap = append(routeOrdersMap, routeOrders{
				RouteReferenceID: inputRoute.ReferenceID,
				Orders:           ordersForRoute,
			})
		}

		// Recopilamos las órdenes sin asignar a ninguna ruta
		for _, order := range p.UnassignedOrders {
			unassignedOrders = append(unassignedOrders, order)
		}

		err = conn.Transaction(func(tx *gorm.DB) error {
			// Inserta o actualiza el plan
			if err := tx.
				Omit("OrganizationCountry").
				Omit("PlanType").
				Omit("PlanningStatus").
				Save(&DBPlanToUpdate).Error; err != nil {
				return err
			}

			// Actualiza las rutas asociadas al plan
			for _, DBRouteToUpdate := range DBRoutesToUpdate {
				if err := tx.
					Omit("OrganizationCountry").
					Omit("Plan").
					Omit("Account").
					Omit("Vehicle").
					Omit("Carrier").
					Save(&DBRouteToUpdate).Error; err != nil {
					return err
				}
			}

			// Reiniciamos la asignación previa en las órdenes
			if err := tx.Table("orders").
				Where("plan_id = ?", DBPlanToUpdate.ID).
				Updates(map[string]interface{}{
					"route_id": nil,
					"plan_id":  nil,
				}).Error; err != nil {
				return err
			}

			// Actualiza individualmente cada orden asignada a una ruta
			for _, ro := range routeOrdersMap {
				if len(ro.Orders) > 0 {
					var route table.Route
					if err := tx.Where("reference_id = ? AND organization_country_id = ?",
						ro.RouteReferenceID, p.Organization.OrganizationCountryID).
						First(&route).Error; err != nil {
						return err
					}
					// Iteramos sobre cada orden individualmente
					for _, order := range ro.Orders {
						planLocation := order.Destination.AddressInfo.PlanLocation
						JSONPlanLocation := table.JSONPlanLocation{
							//NodeReferenceID: string(order.Destination.ReferenceID),
							Longitude: planLocation.Lat(),
							Latitude:  planLocation.Lon(),
						}
						planCorrectedLocation := order.Destination.AddressInfo.PlanCorrectedLocation
						JSONCorrectedPlanLocation := table.JSONPlanLocation{
							//NodeReferenceID: string(order.Destination.ReferenceID),
							Longitude: planCorrectedLocation.Lat(),
							Latitude:  planCorrectedLocation.Lon(),
						}
						correctedDistance := order.Destination.AddressInfo.PlanCorrectedDistance
						if err := tx.Table("orders").
							Where("reference_id = ? AND organization_country_id = ?",
								order.ReferenceID, p.Organization.OrganizationCountryID).
							Updates(map[string]interface{}{
								"route_id":                     route.ID,
								"plan_id":                      DBPlanToUpdate.ID,
								"order_status_id":              plannedStatusID,
								"json_plan_location":           JSONPlanLocation,
								"json_plan_corrected_location": JSONCorrectedPlanLocation,
								"plan_corrected_distance":      correctedDistance,
								"sequence_number":              order.SequenceNumber,
							}).Error; err != nil {
							return err
						}
					}
				}
			}

			// Actualiza individualmente las órdenes sin asignación de ruta
			for _, order := range unassignedOrders {
				planLocation := order.Destination.AddressInfo.PlanLocation
				JSONPlanLocation := table.JSONPlanLocation{
					//NodeReferenceID: string(order.Destination.ReferenceID),
					Longitude: planLocation.Lat(),
					Latitude:  planLocation.Lon(),
				}
				if err := tx.Table("orders").
					Where("reference_id = ? AND organization_country_id = ?",
						order.ReferenceID, p.Organization.OrganizationCountryID).
					Updates(map[string]interface{}{
						"plan_id":            DBPlanToUpdate.ID,
						"order_status_id":    plannedStatusID,
						"json_plan_location": JSONPlanLocation,
					}).Error; err != nil {
					return err
				}
			}

			return nil
		})

		return err
	}
}
