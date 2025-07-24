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

type UpsertContactWorkflow func(ctx context.Context, contact domain.Contact) error

func init() {
	ioc.Registry(
		NewUpsertContactWorkflow,
		workflows.NewUpsertContactWorkflow,
		tidbrepository.NewUpsertContact,
		observability.NewObservability)
}

func NewUpsertContactWorkflow(
	domainWorkflow workflows.UpsertContactWorkflow,
	upsertContact tidbrepository.UpsertContact,
	obs observability.Observability,
) UpsertContactWorkflow {
	return func(ctx context.Context, contact domain.Contact) error {
		key := contact.DocID(ctx).String()
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetContactUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"contact_doc_id", contact.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertContact(ctx, contact, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert contact: %w", err)
		}
		return nil
	}
}
