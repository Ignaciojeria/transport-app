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

type CreateAccountWorkflow func(ctx context.Context, input domain.Account) error

func init() {
	ioc.Registry(
		NewCreateAccountWorkflow,
		workflows.NewCreateAccountWorkflow,
		tidbrepository.NewUpsertAccount,
		observability.NewObservability,
	)
}

func NewCreateAccountWorkflow(
	createAccountWorkflow workflows.CreateAccountWorkflow,
	upsertAccount tidbrepository.UpsertAccount,
	obs observability.Observability) CreateAccountWorkflow {
	return func(ctx context.Context, input domain.Account) error {
		workflow, err := createAccountWorkflow.Restore(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetAccountCreatedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"email", input.Email)
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertAccount(ctx, input, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert account: %w", err)
		}
		return nil
	}
}
