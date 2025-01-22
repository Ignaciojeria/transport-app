package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchNodesQuey func(context.Context, domain.Pagination) ([]domain.NodeInfo, error)

func init() {
	ioc.Registry(NewSearchNodesQuery, tidb.NewTIDBConnection)
}
func NewSearchNodesQuery(conn tidb.TIDBConnection) SearchNodesQuey {
	return func(ctx context.Context, p domain.Pagination) ([]domain.NodeInfo, error) {
		return nil, nil
	}
}
