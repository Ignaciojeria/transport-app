package usecase

import (
	"context"
	"fmt"
	"strings"
	"transport-app/app/shared/infrastructure/cache"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/redis/go-redis/v9"
	"github.com/valkey-io/valkey-go"
)

type GetDataFromRedisWorkflow func(ctx context.Context, key string) ([]byte, error)

func init() {
	ioc.Registry(
		NewGetDataFromRedisWorkflow,
		cache.NewCacheClientFactory,
		observability.NewObservability,
	)
}

func NewGetDataFromRedisWorkflow(
	factory any,
	obs observability.Observability,
) GetDataFromRedisWorkflow {
	return func(ctx context.Context, key string) ([]byte, error) {
		if client, ok := factory.(*redis.Client); ok {
			val, err := client.Get(ctx, key).Bytes()
			if err == redis.Nil {
				obs.Logger.InfoContext(ctx, "Key no encontrada en Redis", "key", key)
				return nil, nil
			}
			if err != nil {
				obs.Logger.ErrorContext(ctx, "Error leyendo desde Redis", "key", key, "error", err)
				return nil, fmt.Errorf("redis get failed: %w", err)
			}
			return val, nil
		}
		if client, ok := factory.(valkey.Client); ok {
			cmd := client.B().Get().Key(key).Build()
			resp := client.Do(ctx, cmd)
			str, err := resp.ToString()
			if err != nil {
				if strings.Contains(err.Error(), "nil") {
					obs.Logger.InfoContext(ctx, "Key no encontrada en Valkey", "key", key)
					return nil, nil
				}
				obs.Logger.ErrorContext(ctx, "Error leyendo desde Valkey", "key", key, "error", err)
				return nil, fmt.Errorf("valkey get failed: %w", err)
			}
			return []byte(str), nil
		}
		return nil, fmt.Errorf("unsupported cache client type")
	}
}
