package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertNode func(context.Context, domain.Origin) error

func init() {
	ioc.Registry(
		NewUpsertNode,
		tidbrepository.NewUpsertNodeQuery,
		tidbrepository.NewUpsertNode,
	)
}

func NewUpsertNode(
	query tidbrepository.UpsertNodeQuery,
	upsert tidbrepository.UpsertNode,
) UpsertNode {
	return func(ctx context.Context, origin domain.Origin) error {
		o, err := query(ctx, origin)
		if err != nil {
			return err
		}
		o.UpdateIfChanged(origin)
		if err := upsert(ctx, o); err != nil {
			return err
		}
		return nil
	}
}
