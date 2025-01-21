package tidbrepository

import (
	"context"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertVehicle func(context.Context, domain.Vehicle) error

func init() {
	ioc.Registry(NewUpsertVehicle, tidb.NewTIDBConnection)
}
func NewUpsertVehicle(conn tidb.TIDBConnection) UpsertVehicle {
	return func(ctx context.Context, v domain.Vehicle) error {
		tbl := mapper.DomainToTableVehicle(v)
		tbl.CreatedAt = time.Now()
		tbl.OrganizationCountryID = v.Organization.OrganizationCountryID
		tbl.Carrier.OrganizationCountryID = v.Organization.OrganizationCountryID
		tbl.Carrier.CreatedAt = time.Now()
		if err := conn.Save(&tbl.Carrier).Error; err != nil {
			return err
		}
		return conn.Save(&tbl).Error
	}
}
