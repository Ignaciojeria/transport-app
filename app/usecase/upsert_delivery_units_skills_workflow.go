package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertDeliveryUnitsSkillsWorkflow func(ctx context.Context, order domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsSkillsWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertDeliveryUnitsSkills,
		observability.NewObservability)
}

func NewUpsertDeliveryUnitsSkillsWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertDeliveryUnitsSkills tidbrepository.UpsertDeliveryUnitsSkills,
	obs observability.Observability,
) UpsertDeliveryUnitsSkillsWorkflow {
	return func(ctx context.Context, order domain.Order) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("delivery_units_skills")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"delivery_units_skills_doc_id", order.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertDeliveryUnitsSkills(ctx, order, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert delivery units skills: %w", err)
		}
		return nil
	}
}