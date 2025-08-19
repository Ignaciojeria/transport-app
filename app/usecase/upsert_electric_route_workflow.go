package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/domain/workflows"
	canonicaljson "transport-app/app/shared/caonincaljson"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertElectricRouteWorkflow func(ctx context.Context, route domain.Route, planDoc string, input interface{}) error

func init() {
	ioc.Registry(NewUpsertElectricRouteWorkflow,
		workflows.NewUpsertElectricRouteWorkflow,
		tidbrepository.NewUpsertRoute,
		observability.NewObservability)
}

func NewUpsertElectricRouteWorkflow(
	workflow workflows.UpsertElectricRouteWorkflow,
	upsertRoute tidbrepository.UpsertRoute,
	obs observability.Observability,
) UpsertElectricRouteWorkflow {
	return func(ctx context.Context, route domain.Route, planDoc string, input interface{}) error {
		// Generar clave de idempotencia para el workflow
		key, err := canonicaljson.HashKey(ctx, "electric_route", map[string]interface{}{
			"route":   route,
			"planDoc": planDoc,
			"input":   input,
		})
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}

		// Restaurar el workflow usando la clave de idempotencia
		workflowInstance, err := workflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}

		// Intentar hacer la transición de estado
		if err := workflowInstance.SetElectricRouteUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"route_doc_id", route.DocID(ctx).String(),
				"plan_doc", planDoc)
			return nil
		}

		// Mapear el estado del workflow
		fsmState := workflowInstance.Map(ctx)

		// Ejecutar el upsert de la ruta con el estado del workflow
		err = upsertRoute(ctx, route, input, planDoc, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert electric route: %w", err)
		}

		return nil
	}
}
