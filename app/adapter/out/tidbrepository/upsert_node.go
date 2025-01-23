package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertNode func(context.Context, domain.NodeInfo) error

func init() {
	ioc.Registry(NewUpsertNode, tidb.NewTIDBConnection)
}
func NewUpsertNode(conn tidb.TIDBConnection) UpsertNode {
	return func(ctx context.Context, ni domain.NodeInfo) error {
		return nil
	}
}
