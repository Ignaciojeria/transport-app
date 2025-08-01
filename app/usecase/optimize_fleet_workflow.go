package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/adapter/out/vroom"
	"transport-app/app/domain/optimization"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OptimizeFleetWorkflow func(ctx context.Context, input optimization.FleetOptimization) ([]request.UpsertRouteRequest, error)

func init() {
	ioc.Registry(
		NewOptimizeFleetWorkflow,
		workflows.NewOptimizeFleetWorkflow,
		vroom.NewOptimize,
		storjbucket.NewTransportAppBucket,
		tidbrepository.NewSaveFSMTransition,
		observability.NewObservability,
	)
}

func NewOptimizeFleetWorkflow(
	domainWorkflow workflows.OptimizeFleetWorkflow,
	optimize vroom.Optimize,
	storjBucket *storjbucket.TransportAppBucket,
	saveFSMTransition tidbrepository.SaveFSMTransition,
	obs observability.Observability,
) OptimizeFleetWorkflow {
	return func(ctx context.Context, input optimization.FleetOptimization) ([]request.UpsertRouteRequest, error) {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetOptimizationCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error())
			fsmState := workflow.Map(ctx)
			var optimizeFleetWorkflowNextInput []request.UpsertRouteRequest
			if err := json.Unmarshal(fsmState.NextInput, &optimizeFleetWorkflowNextInput); err != nil {
				obs.Logger.ErrorContext(ctx, "Error deserializando payload de optimización (reconstruido)", "error", err)
				return nil, fmt.Errorf("error deserializing optimization payload: %w", err)
			}
			return optimizeFleetWorkflowNextInput, nil
		}
		routeRequests, err := optimize(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to optimize fleet: %w", err)
		}

		// Serializar correctamente las rutas para NextInput
		routeRequestsJSON, err := json.Marshal(routeRequests)
		if err != nil {
			return nil, fmt.Errorf("error marshaling route requests: %w", err)
		}

		fsmState := workflow.Map(ctx)
		fsmState.NextInput = routeRequestsJSON
		err = saveFSMTransition(ctx, fsmState)
		if err != nil {
			return nil, fmt.Errorf("failed to save FSM transition: %w", err)
		}
		fmt.Printf("Optimización completada. Se generaron %d rutas.\n", len(routeRequests))
		return routeRequests, nil
	}
}
