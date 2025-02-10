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

		routeMap := make(map[string]table.Route)
		for _, r := range routes {
			routeMap[r.ReferenceID] = r
		}

		var DBRoutesToUpdate []table.Route
		type routeOrders struct {
			RouteID  string
			OrderIDs []string
		}
		var routeOrdersMap []routeOrders
		var unassignedOrderRefs []string

		for _, inputRoute := range p.Routes {
			existingRoute, exists := routeMap[inputRoute.ReferenceID]

			var routeToUpdate domain.Route
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

			var orderRefs []string
			for _, order := range inputRoute.Orders {
				orderRefs = append(orderRefs, string(order.ReferenceID))
			}

			routeOrdersMap = append(routeOrdersMap, routeOrders{
				RouteID:  inputRoute.ReferenceID,
				OrderIDs: orderRefs,
			})
		}

		for _, order := range p.UnassignedOrders {
			unassignedOrderRefs = append(unassignedOrderRefs, string(order.ReferenceID))
		}

		err = conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.
				Omit("OrganizationCountry").
				Omit("PlanType").
				Omit("PlanningStatus").
				Save(&DBPlanToUpdate).Error; err != nil {
				return err
			}

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

			if err := tx.Table("orders").
				Where("plan_id = ?", DBPlanToUpdate.ID).
				Updates(map[string]interface{}{
					"route_id": nil,
					"plan_id":  nil,
				}).Error; err != nil {
				return err
			}

			plannedStatusID := loadOrdorderStatuses().Planned().ID

			for _, ro := range routeOrdersMap {
				if len(ro.OrderIDs) > 0 {
					var route table.Route
					if err := tx.Where("reference_id = ?", ro.RouteID).First(&route).Error; err != nil {
						return err
					}

					if err := tx.Table("orders").
						Where("reference_id IN ? AND organization_country_id = ?",
							ro.OrderIDs, p.Organization.OrganizationCountryID).
						Updates(map[string]interface{}{
							"route_id":        route.ID,
							"plan_id":         DBPlanToUpdate.ID,
							"order_status_id": plannedStatusID,
						}).Error; err != nil {
						return err
					}
				}
			}

			if len(unassignedOrderRefs) > 0 {
				if err := tx.Table("orders").
					Where("reference_id IN ? AND organization_country_id = ?",
						unassignedOrderRefs, p.Organization.OrganizationCountryID).
					Updates(map[string]interface{}{
						"plan_id":         DBPlanToUpdate.ID,
						"order_status_id": plannedStatusID,
					}).Error; err != nil {
					return err
				}
			}

			return nil
		})

		return err
	}
}
