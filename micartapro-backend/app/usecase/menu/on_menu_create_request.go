package menu

import (
	"context"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnMenuCreateRequest func(ctx context.Context, input events.MenuCreateRequest) error

func init() {
	ioc.Registry(NewOnMenuCreateRequest,
		observability.NewObservability,
		storage.NewSaveMenu,
	)
}

func NewOnMenuCreateRequest(
	observability observability.Observability,
	saveMenu storage.SaveMenu) OnMenuCreateRequest {
	return func(ctx context.Context, input events.MenuCreateRequest) error {
		observability.Logger.Info("on_menu_create_request", "input", input)
		spanCtx, span := observability.Tracer.Start(ctx, "on_menu_create_request")
		defer span.End()
		err := saveMenu(spanCtx, input)
		if err != nil {
			observability.Logger.Error("error_saving_menu", "error", err)
			return err
		}
		return nil
	}
}
