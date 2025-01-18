package tidbrepository

import (
	"context"
	"transport-app/app/domain"
)

type UpsertNodeQuery func(ctx context.Context, origin domain.Origin)
