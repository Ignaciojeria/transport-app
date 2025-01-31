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

type UpsertDailyPlan func(context.Context, domain.Plan) error

func init() {
	ioc.Registry(
		NewUpsertDailyPlan,
		tidb.NewTIDBConnection)
}

func NewUpsertDailyPlan(conn tidb.TIDBConnection) UpsertDailyPlan {
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

		route := table.Route{}
		err = conn.DB.WithContext(ctx).Table("routes").
			Where("reference_id = ? AND organization_country_id = ?",
				p.ReferenceID, p.Organization.OrganizationCountryID).First(&route).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		p.Routes[0].PlanID = DBPlanToUpdate.ID
		routeToUpdate := route.Map().UpdateIfChanged(p.Routes[0])
		DBRouteToUpdate := mapper.MapRoute(routeToUpdate)
		DBRouteToUpdate.CreatedAt = route.CreatedAt

		orderList := []int64{}

		for _, order := range p.Routes[0].Orders {
			orderList = append(orderList, order.ID)
		}

		err = conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.
				Omit("OrganizationCountry").
				Omit("PlanType").
				Omit("PlanningStatus").
				Save(&DBPlanToUpdate).Error; err != nil {
				return err
			}
			if err := tx.
				Omit("OrganizationCountry").
				Omit("Plan").
				Omit("Account").
				Omit("Vehicle").
				Omit("Carrier").
				Save(&DBRouteToUpdate).Error; err != nil {
				return err
			}

			if err := tx.Table("orders").
				Where("route_id = ?", DBRouteToUpdate.ID).
				Update("route_id", nil).Error; err != nil {
				return err
			}

			if len(orderList) > 0 {
				if err := tx.Table("orders").
					Where("id IN ?", orderList).
					Update("route_id", DBRouteToUpdate.ID).Error; err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}
}
