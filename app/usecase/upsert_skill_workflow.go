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

type UpsertSkillWorkflow func(ctx context.Context, skill domain.Skill) error

func init() {
	ioc.Registry(
		NewUpsertSkillWorkflow,
		workflows.NewUpsertSkillWorkflow,
		tidbrepository.NewUpsertSkill,
		observability.NewObservability)
}

func NewUpsertSkillWorkflow(
	domainWorkflow workflows.UpsertSkillWorkflow,
	upsertSkill tidbrepository.UpsertSkill,
	obs observability.Observability,
) UpsertSkillWorkflow {
	return func(ctx context.Context, skill domain.Skill) error {
		key, err := canonicaljson.HashKey(ctx, "skill", skill)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetSkillUpsertedTransition(ctx); err != nil {
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
