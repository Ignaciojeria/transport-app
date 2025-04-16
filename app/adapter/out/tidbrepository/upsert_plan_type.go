package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertPlanType func(context.Context, domain.PlanType) error

func init() {
	ioc.Registry(NewUpsertPlanType, database.NewConnectionFactory)
}
func NewUpsertPlanType(conn database.ConnectionFactory) UpsertPlanType {
	return func(ctx context.Context, pt domain.PlanType) error {
		var planType table.PlanType
		err := conn.DB.WithContext(ctx).
			Table("plan_types").
			Where("document_id = ?", pt.DocID(ctx)).
			First(&planType).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return conn.
				Omit("Organization").
				Save(mapper.MapPlanType(ctx, pt)).Error
		}
		return nil
	}
}
