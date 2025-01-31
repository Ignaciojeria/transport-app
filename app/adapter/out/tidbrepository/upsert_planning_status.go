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

type UpsertPlanningStatus func(context.Context, domain.PlanningStatus) (domain.PlanningStatus, error)

func init() {
	ioc.Registry(NewUpsertPlanningStatus, tidb.NewTIDBConnection)
}
func NewUpsertPlanningStatus(conn tidb.TIDBConnection) UpsertPlanningStatus {
	return func(ctx context.Context, ps domain.PlanningStatus) (domain.PlanningStatus, error) {
		var planningStatus table.PlanningStatus
		err := conn.DB.WithContext(ctx).
			Table("planning_statuses").
			Where("name = ? AND organization_country_id = ?", ps.Value, ps.Organization.OrganizationCountryID).
			First(&planningStatus).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.PlanningStatus{}, err
		}
		planningStatusWithChanges := planningStatus.Map().UpdateIfChanged(ps)
		DBPlanningStatusToUpdate := mapper.MapPlanningStatus(planningStatusWithChanges)
		DBPlanningStatusToUpdate.CreatedAt = planningStatus.CreatedAt
		if err := conn.
			Omit("OrganizationCountry").
			Save(&DBPlanningStatusToUpdate).Error; err != nil {
			return domain.PlanningStatus{}, err
		}
		return DBPlanningStatusToUpdate.Map(), nil
	}
}
