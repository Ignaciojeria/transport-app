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

type UpsertSizeCategory func(context.Context, domain.SizeCategory) error

func init() {
	ioc.Registry(NewUpsertSizeCategory, database.NewConnectionFactory)
}

func NewUpsertSizeCategory(conn database.ConnectionFactory) UpsertSizeCategory {
	return func(ctx context.Context, sc domain.SizeCategory) error {
		var existing table.SizeCategory

		err := conn.DB.WithContext(ctx).
			Table("size_categories").
			Where("document_id = ?", sc.DocumentID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			// Ya existe → no hacer nada
			return nil
		}

		newRecord := mapper.MapSizeCategory(ctx, sc)

		err = conn.DB.WithContext(ctx).Create(&newRecord).Error
		if err != nil {
			return err
		}

		return nil
	}
}
