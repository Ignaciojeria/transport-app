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

type UpsertPlanningStatus func(context.Context, domain.PlanningStatus) error

func init() {
	ioc.Registry(NewUpsertPlanningStatus, tidb.NewTIDBConnection)
}

func NewUpsertPlanningStatus(conn tidb.TIDBConnection) UpsertPlanningStatus {
	return func(ctx context.Context, ps domain.PlanningStatus) error {
		var planningStatus table.PlanningStatus
		err := conn.DB.WithContext(ctx).
			Table("planning_statuses").
			Where("document_id = ?", ps.Value, ps.Organization.ID).
			First(&planningStatus).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return conn.
				Omit("Organization").
				Save(mapper.MapPlanningStatus(ps)).Error
		}
		return nil
	}
}
