package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type StoreDataInBucketWorkflow func(ctx context.Context, key string, value []byte) error

func init() {
	ioc.Registry(NewStoreDataInBucketWorkflow,
		workflows.NewStoreDataInBucketWorkflow,
		storjbucket.NewTransportAppBucket,
		tidbrepository.NewSaveFSMTransition,
		observability.NewObservability,
	)
}

func NewStoreDataInBucketWorkflow(
	domainWorkflow workflows.StoreDataInBucketWorkflow,
	storjBucket *storjbucket.TransportAppBucket,
	saveFSMTransition tidbrepository.SaveFSMTransition,
	obs observability.Observability,
) StoreDataInBucketWorkflow {
	return func(ctx context.Context, key string, value []byte) error {
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return err
		}
		if err := workflow.SetDataStoredInBucketTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error())
			return nil
		}
		token, ok := sharedcontext.BucketTokenFromContext(ctx)
		if !ok {
			return fmt.Errorf("bucket token not found in context")
		}
		err = storjBucket.UploadWithToken(ctx, token, key, value)
		if err != nil {
			return fmt.Errorf("error uploading route request: %w", err)
		}
		fsmState := workflow.Map(ctx)
		err = saveFSMTransition(ctx, fsmState)
		if err != nil {
			return fmt.Errorf("failed to save FSM transition: %w", err)
		}
		return nil
	}
}
