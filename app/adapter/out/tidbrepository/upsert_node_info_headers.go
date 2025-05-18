package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertNodeInfoHeaders func(context.Context, domain.Headers) error

func init() {
	ioc.Registry(NewUpsertNodeInfoHeaders, database.NewConnectionFactory)
}

func NewUpsertNodeInfoHeaders(conn database.ConnectionFactory) UpsertNodeInfoHeaders {
	return func(ctx context.Context, h domain.Headers) error {

		nodeInfoHeadersTbl := mapper.MapNodeInfoHeaders(ctx, h)

		err := conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("document_id = ?", h.DocID(ctx)).
			First(&nodeInfoHeadersTbl).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			// Ya existe → no hacer nada
			return nil
		}

		return conn.DB.
			WithContext(ctx).
			Omit("Tenant").
			Create(&nodeInfoHeadersTbl).Error
	}
}
