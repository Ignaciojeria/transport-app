package cacherepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/cache"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/paulmach/orb"
	"github.com/redis/go-redis/v9"
	"github.com/valkey-io/valkey-go"
)

type GeocodingCacheStrategy interface {
	Save(context.Context, domain.AddressInfo) error
	Get(context.Context, domain.AddressInfo) (orb.Point, error)
}

func init() {
	ioc.Registry(NewGeocodingCacheStrategy, cache.NewCacheClientFactory)
}
func NewGeocodingCacheStrategy(factory any) GeocodingCacheStrategy {
	if client, ok := factory.(*redis.Client); ok {
		return newRedisGeocodingCacheStrategy(client)
	}
	if client, ok := factory.(valkey.Client); ok {
		return newValkeyGeocodingCacheStrategy(client)
	}
	panic("unsupported cache client type")
}
