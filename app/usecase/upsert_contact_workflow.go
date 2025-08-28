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

type UpsertContactWorkflow func(ctx context.Context, contact domain.Contact) error

func init() {
	ioc.Registry(
		NewUpsertContactWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertContact,
		observability.NewObservability)
}

func NewUpsertContactWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertContact tidbrepository.UpsertContact,
	obs observability.Observability,
) UpsertContactWorkflow {
	return func(ctx context.Context, contact domain.Contact) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("contact")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
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
