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

type UpsertProvince func(ctx context.Context, province domain.Province) error

func init() {
	ioc.Registry(NewUpsertProvince, database.NewConnectionFactory)
}

func NewUpsertProvince(db database.ConnectionFactory) UpsertProvince {
	return func(ctx context.Context, province domain.Province) error {
		var existing table.Province
		err := db.WithContext(ctx).
			Table("provinces").
			Where("document_id = ?", province.DocID(ctx).String()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Si ya existe, no hacemos nada porque el nombre es el mismo
		if err == nil {
			return nil
		}

		// No existe → insert
		newProvince := mapper.MapProvinceTable(ctx, province)
		return db.Create(&newProvince).Error
	}
}
