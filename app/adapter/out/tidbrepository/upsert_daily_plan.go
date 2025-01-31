package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	"gorm.io/gorm"
)

type UpsertDailyPlan func(context.Context, domain.Plan) (domain.Plan, error)

func NewUpsertDailyPlan(conn tidb.TIDBConnection) UpsertDailyPlan {
	return func(ctx context.Context, p domain.Plan) (domain.Plan, error) {
		plan := table.Plan{}
		err := conn.DB.WithContext(ctx).Table("plans").
			Where("reference_id = ? AND organization_country_id = ?",
				p.ReferenceID, p.Organization.OrganizationCountryID).First(&plan).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Plan{}, err
		}
		planToUpdate := plan.Map().UpdateIfChanged(p)
		DBPlanToUpdate := mapper.MapPlan(planToUpdate)
		DBPlanToUpdate.CreatedAt = plan.CreatedAt
		err = conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.
				Omit("OrganizationCountry").
				Omit("PlanType").
				Omit("PlanningStatus").
				Save(&DBPlanToUpdate).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return domain.Plan{}, err
		}

		return domain.Plan{}, nil
	}
}
