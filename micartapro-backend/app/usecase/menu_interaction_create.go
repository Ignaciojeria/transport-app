package usecase

import (
	"context"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/domain"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type MenuInteractionCreate func(ctx context.Context, input domain.MenuCreateRequest) error

func init() {
	ioc.Registry(NewMenuInteractionCreate,
		observability.NewObservability,
		storage.NewSaveMenu,
	)
}

func NewMenuInteractionCreate(
	observability observability.Observability,
	saveMenu storage.SaveMenu) MenuInteractionCreate {
	return func(ctx context.Context, input domain.MenuCreateRequest) error {
		observability.Logger.Info("menu_interaction_create", "input", input)
		spanCtx, span := observability.Tracer.Start(ctx, "menu_interaction_create")
		defer span.End()
		err := saveMenu(spanCtx, input)
		if err != nil {
			observability.Logger.Error("error_saving_menu", "error", err)
			return err
		}
		return nil
	}
}
