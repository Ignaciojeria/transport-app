package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
)

type UpsertWebhookWorkflow func(ctx context.Context, w domain.Webhook) error

func init() {
	ioc.Registry(
		NewUpsertWebhookWorkflow,
		tidbrepository.NewUpsertWebhook,
		observability.NewObservability,
		natsconn.NewKeyValue)
}

func NewUpsertWebhookWorkflow(
	upsertWebhook tidbrepository.UpsertWebhook,
	obs observability.Observability,
	kv jetstream.KeyValue,
) UpsertWebhookWorkflow {
	return func(ctx context.Context, w domain.Webhook) error {
		bytes, err := json.Marshal(w)
		if err != nil {
			return fmt.Errorf("failed to marshal webhook: %w", err)
		}
		kv.Put(ctx, w.DocID(ctx).String(), bytes)
		return nil
	}
}
