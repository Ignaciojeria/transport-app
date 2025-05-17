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

type UpsertDistrict func(ctx context.Context, district domain.District) error

func init() {
	ioc.Registry(NewUpsertDistrict, database.NewConnectionFactory)
}

func NewUpsertDistrict(db database.ConnectionFactory) UpsertDistrict {
	return func(ctx context.Context, district domain.District) error {
		var existing table.District
		err := db.WithContext(ctx).
			Table("districts").
			Where("document_id = ?", district.DocID(ctx).String()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Si ya existe, no hacemos nada porque el nombre es el mismo
		if err == nil {
			return nil
		}

		// No existe → insert
		newDistrict := mapper.MapDistrictTable(ctx, district)
		return db.Create(&newDistrict).Error
	}
}
