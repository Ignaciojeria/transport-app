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

type UpsertSkillWorkflow func(ctx context.Context, skill domain.Skill) error

func init() {
	ioc.Registry(
		NewUpsertSkillWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertSkill,
		observability.NewObservability)
}

func NewUpsertSkillWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertSkill tidbrepository.UpsertSkill,
	obs observability.Observability,
) UpsertSkillWorkflow {
	return func(ctx context.Context, skill domain.Skill) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("skill")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"skill_doc_id", key)
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertSkill(ctx, skill, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert skill: %w", err)
		}
		return nil
	}
}
