package menu

import (
	"context"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnMenuCreateRequest func(ctx context.Context, input events.MenuCreateRequest) error

func init() {
	ioc.Registry(NewOnMenuCreateRequest,
		observability.NewObservability,
		storage.NewSaveMenu,
		supabaserepo.NewSaveMenu,
	)
}

func NewOnMenuCreateRequest(
	observability observability.Observability,
	saveMenuStorage storage.SaveMenu,
	saveMenuSupabase supabaserepo.SaveMenu) OnMenuCreateRequest {
	return func(ctx context.Context, input events.MenuCreateRequest) error {
		observability.Logger.Info("on_menu_create_request", "input", input)
		spanCtx, span := observability.Tracer.Start(ctx, "on_menu_create_request")
		defer span.End()

		// Guardar en GCS (storage)
		err := saveMenuStorage(spanCtx, input)
		if err != nil {
			observability.Logger.Error("error_saving_menu_to_storage", "error", err)
			return err
		}

		// Guardar en Supabase
		err = saveMenuSupabase(spanCtx, input)
		if err != nil {
			observability.Logger.Error("error_saving_menu_to_supabase", "error", err)
			return err
		}

		return nil
	}
}
