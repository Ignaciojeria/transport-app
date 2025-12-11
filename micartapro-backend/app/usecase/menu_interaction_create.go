package usecase

import (
	"context"
	"micartapro/app/domain"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type MenuInteractionCreate func(ctx context.Context, input domain.MenuCreateRequest) error

func init() {
	ioc.Registry(NewMenuInteractionCreate,
		observability.NewObservability)
}

func NewMenuInteractionCreate(
	observability observability.Observability) MenuInteractionCreate {
	return func(ctx context.Context, input domain.MenuCreateRequest) error {
		observability.Logger.Info("menu_interaction_create", "input", input)
		return nil
	}
}
