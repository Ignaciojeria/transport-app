package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchNodes func(ctx context.Context, input domain.NodeSearchFilters) ([]domain.NodeInfo, error)

func init() {
	ioc.Registry(NewSearchNodes, tidbrepository.NewSearchNodesQuery)
}

func NewSearchNodes(search tidbrepository.SearchNodesQuery) SearchNodes {
	return func(ctx context.Context, input domain.NodeSearchFilters) ([]domain.NodeInfo, error) {
		return search(ctx, input)
	}
}
