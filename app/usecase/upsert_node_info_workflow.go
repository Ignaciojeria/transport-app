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

type UpsertNodeInfoWorkflow func(ctx context.Context, ni domain.NodeInfo) error

func init() {
	ioc.Registry(
		NewUpsertNodeInfoWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertNodeInfo,
		observability.NewObservability)
}

func NewUpsertNodeInfoWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
	obs observability.Observability,
) UpsertNodeInfoWorkflow {
	return func(ctx context.Context, ni domain.NodeInfo) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "node_info", ni)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		config := workflows.CreateUpsertWorkflow("node_info")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"node_info_doc_id", ni.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertNodeInfo(ctx, ni, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert node info: %w", err)
		}
		return nil
	}
}