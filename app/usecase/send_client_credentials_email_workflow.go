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
		workflows.NewGenericWorkflow,
		tidbrepository.NewFindClientCredentialsByClientID,
		email.NewSendClientCredentialsEmail,
		observability.NewObservability,
		tidbrepository.NewSaveFSMTransition,
	)
}

func NewSendClientCredentialsEmailWorkflow(
	genericWorkflow workflows.GenericWorkflow,
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
		
		// Configurar workflow genérico para email sending
		config := workflows.WorkflowConfig{
			Name:           "send_client_credentials_email_workflow",
			StartedState:   "email_sending_started",
			CompletedState: "email_sent",
			UseStorjBucket: false,
		}
		
		w, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}

		// Intentar transición a email enviado
		if err := w.SetCompletedTransition(ctx); err != nil {
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
