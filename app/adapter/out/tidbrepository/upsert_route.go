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

type UpsertRoute func(context.Context, domain.Route, string) error

func init() {
	ioc.Registry(NewUpsertRoute, database.NewConnectionFactory)
}

func NewUpsertRoute(conn database.ConnectionFactory) UpsertRoute {
	return func(ctx context.Context, r domain.Route, planDoc string) error {
		var existing table.Route

		err := conn.DB.WithContext(ctx).
			Table("routes").
			Where("document_id = ?", r.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newRoute := mapper.MapRouteTable(ctx, r, planDoc)
			return conn.Create(&newRoute).Error
		}

		// Ya existe → update
		updateData := mapper.MapRouteTable(ctx, r, planDoc)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Save(&updateData).Error
	}
}
