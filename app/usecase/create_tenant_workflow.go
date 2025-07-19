package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateTenantWorkflow func(ctx context.Context, input domain.Tenant) error

func init() {
	ioc.Registry(
		NewCreateTenantWorkflow,
		workflows.NewCreateTenantWorkflow,
		tidbrepository.NewSaveTenant,
		observability.NewObservability,
	)
}

func NewCreateTenantWorkflow(
	createTenantWorkflow workflows.CreateTenantWorkflow,
	saveTenant tidbrepository.SaveTenant,
	obs observability.Observability) CreateTenantWorkflow {
	return func(ctx context.Context, input domain.Tenant) error {
		workflow, err := createTenantWorkflow.Restore(ctx, input.ID.String())
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetTenantCreatedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"tenant_id", input.ID.String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		_, err = saveTenant(ctx, input, fsmState)
		if err != nil {
			return fmt.Errorf("failed to save tenant: %w", err)
		}
		return nil
	}
}
