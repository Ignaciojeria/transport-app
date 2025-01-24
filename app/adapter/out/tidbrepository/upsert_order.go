package tidbrepository

import (
	"context"
	"fmt"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertOrder func(context.Context, domain.Order) (domain.Order, error)

func init() {
	ioc.Registry(NewUpsertOrder, tidb.NewTIDBConnection)
}
func NewUpsertOrder(conn tidb.TIDBConnection) UpsertOrder {
	return func(ctx context.Context, o domain.Order) (domain.Order, error) {
		fmt.Println(o)
		return domain.Order{}, nil
	}
}
