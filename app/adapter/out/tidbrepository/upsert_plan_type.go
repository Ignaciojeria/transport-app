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

type UpsertPlanType func(context.Context, domain.PlanType) (domain.PlanType, error)

func init() {
	ioc.Registry(NewUpsertPlanType, tidb.NewTIDBConnection)
}
func NewUpsertPlanType(conn tidb.TIDBConnection) UpsertPlanType {
	return func(ctx context.Context, pt domain.PlanType) (domain.PlanType, error) {
		var planType table.PlanType
		err := conn.DB.WithContext(ctx).
			Table("plan_types").
			Where("name = ? AND organization_country_id = ?", pt.Value, pt.Organization.OrganizationCountryID).
			First(&planType).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.PlanType{}, err
		}
		planTypeWithChanges := planType.Map().UpdateIfChanged(pt)
		DBOrderTypeToUpdate := mapper.MapPlanType(planTypeWithChanges)
		DBOrderTypeToUpdate.CreatedAt = planType.CreatedAt
		if err := conn.
			Omit("OrganizationCountry").
			Save(&DBOrderTypeToUpdate).Error; err != nil {
			return domain.PlanType{}, err
		}
		return DBOrderTypeToUpdate.Map(), nil
	}
}
