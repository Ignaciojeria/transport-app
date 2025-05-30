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

type UpsertNodeType func(context.Context, domain.NodeType) error

func init() {
	ioc.Registry(NewUpsertNodeType, database.NewConnectionFactory)
}

func NewUpsertNodeType(conn database.ConnectionFactory) UpsertNodeType {
	return func(ctx context.Context, nt domain.NodeType) error {
		var existing table.NodeType

		err := conn.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nt.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			// Ya existe → no hacer nada
			return nil
		}
		newRecord := mapper.MapNodeType(ctx, nt)

		err = conn.Omit("Tenant").Create(&newRecord).Error
		if err != nil {
			return err
		}

		return nil
	}
}
