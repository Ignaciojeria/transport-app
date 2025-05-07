package tidbrepository

import (
	"context"
	"errors"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/paulmach/orb"
	"gorm.io/gorm"
)

type GetStoredCoordinates func(ctx context.Context, adi domain.AddressInfo) (orb.Point, error)

func init() {
	ioc.Registry(NewGetStoredCoordinates, database.NewConnectionFactory)
}

func NewGetStoredCoordinates(conn database.ConnectionFactory) GetStoredCoordinates {
	return func(ctx context.Context, adi domain.AddressInfo) (orb.Point, error) {
		var row table.AddressInfo
		err := conn.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", adi.DocID(ctx)).
			Select("latitude", "longitude").
			First(&row).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return orb.Point{}, nil // No hay coordenadas aún
		}
		if err != nil {
			return orb.Point{}, fmt.Errorf("db error getting coordinates: %w", err)
		}

		return orb.Point{row.Longitude, row.Latitude}, nil
	}
}
