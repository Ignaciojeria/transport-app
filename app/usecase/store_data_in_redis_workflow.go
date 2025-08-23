package usecase

import (
	"context"
	"fmt"
	"transport-app/app/shared/infrastructure/cache"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/redis/go-redis/v9"
	"github.com/valkey-io/valkey-go"
)

type StoreDataInRedisWorkflow func(ctx context.Context, key string, value []byte) error

func init() {
	ioc.Registry(
		NewStoreDataInRedisWorkflow,
		cache.NewCacheClientFactory,
		observability.NewObservability,
	)
}

func NewStoreDataInRedisWorkflow(
	factory any,
	obs observability.Observability,
) StoreDataInRedisWorkflow {
	return func(ctx context.Context, key string, value []byte) error {
		if client, ok := factory.(*redis.Client); ok {
			if err := client.Set(ctx, key, value, 0).Err(); err != nil {
				obs.Logger.ErrorContext(ctx, "Error guardando en Redis", "key", key, "error", err)
				return fmt.Errorf("redis set failed: %w", err)
			}
			obs.Logger.InfoContext(ctx, "Dato guardado en Redis", "key", key)
			return nil
		}
		if client, ok := factory.(valkey.Client); ok {
			cmd := client.B().Set().Key(key).Value(string(value)).Build()
			if err := client.Do(ctx, cmd).Error(); err != nil {
				obs.Logger.ErrorContext(ctx, "Error guardando en Valkey", "key", key, "error", err)
				return fmt.Errorf("valkey set failed: %w", err)
			}
			obs.Logger.InfoContext(ctx, "Dato guardado en Valkey", "key", key)
			return nil
		}
		return fmt.Errorf("unsupported cache client type")
	}
}
