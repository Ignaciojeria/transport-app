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

type UpsertPlanType func(context.Context, domain.PlanType) error

func init() {
	ioc.Registry(NewUpsertPlanType, tidb.NewTIDBConnection)
}
func NewUpsertPlanType(conn tidb.TIDBConnection) UpsertPlanType {
	return func(ctx context.Context, pt domain.PlanType) error {
		var planType table.PlanType
		err := conn.DB.WithContext(ctx).
			Table("plan_types").
			Where("document_id = ?", pt.DocID()).
			First(&planType).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return conn.
				Omit("Organization").
				Save(mapper.MapPlanType(pt)).Error
		}
		return nil
	}
}
