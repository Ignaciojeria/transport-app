package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertNodeQuery func(
	ctx context.Context,
	origin domain.Origin) (domain.Origin, error)

func init() {
	ioc.Registry(
		NewUpsertNodeQuery,
		tidb.NewTIDBConnection)
}
func NewUpsertNodeQuery(conn tidb.TIDBConnection) UpsertNodeQuery {
	return func(ctx context.Context, origin domain.Origin) (domain.Origin, error) {

		return domain.Origin{}, nil
	}
}
