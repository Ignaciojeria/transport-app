package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/email"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SendClientCredentialsEmailWorkflow func(ctx context.Context, clientID string, email string) error

func init() {
	ioc.Registry(
		NewSendClientCredentialsEmailWorkflow,
		workflows.NewSendClientCredentialsEmailWorkflow,
		tidbrepository.NewFindClientCredentialsByClientID,
		email.NewSendClientCredentialsEmail,
		observability.NewObservability,
		tidbrepository.NewSaveFSMTransition,
	)
}

func NewSendClientCredentialsEmailWorkflow(
	workflow workflows.SendClientCredentialsEmailWorkflow,
	findClientCredentials tidbrepository.FindClientCredentialsByClientID,
	sendEmail email.SendClientCredentialsEmail,
	obs observability.Observability,
	saveFSMTransition tidbrepository.SaveFSMTransition,
) SendClientCredentialsEmailWorkflow {
	return func(ctx context.Context, clientID string, email string) error {
		// Obtener el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		w, err := workflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}

		// Intentar transición a email enviado
		if err := w.SetEmailSentTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"client_id", clientID)
			return nil
		}

		// Buscar las credenciales por client ID
		credentials, err := findClientCredentials(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to find client credentials: %w", err)
		}
		// Enviar el email con las credenciales
		err = sendEmail(ctx, email, credentials)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		fsmState := w.Map(ctx)
		err = saveFSMTransition(ctx, fsmState)
		if err != nil {
			return fmt.Errorf("failed to save FSM transition: %w", err)
		}

		return nil
	}
}
