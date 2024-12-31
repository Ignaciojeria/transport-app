package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewSaveTransportOrder, tidb.NewTIDBConnection)
}

type SaveTransportOrder func(context.Context, domain.TransportOrder) (domain.TransportOrder, error)

func NewSaveTransportOrder(conn tidb.TIDBConnection) SaveTransportOrder {
	return func(ctx context.Context, to domain.TransportOrder) (domain.TransportOrder, error) {

		return domain.TransportOrder{}, nil
	}
}
