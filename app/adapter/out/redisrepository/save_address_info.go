package redisrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/rediscache"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type SaveAddressInfo func(ctx context.Context, raw, normalized domain.AddressInfo) error

func init() {
	ioc.Registry(NewSaveAddressInfo, rediscache.NewRedisClient)
}

func NewSaveAddressInfo(c *redis.Client) SaveAddressInfo {
	return func(ctx context.Context, raw, normalized domain.AddressInfo) error {
		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			if raw.Province.Equals(normalized.Province) {
				return nil
			}
			return c.Set(ctx, raw.Province.DocID(ctx).String(), normalized.Province.String(), 0).Err()
		})

		group.Go(func() error {
			if raw.State.Equals(normalized.State) {
				return nil
			}
			return c.Set(ctx, raw.State.DocID(ctx).String(), normalized.State.String(), 0).Err()
		})

		group.Go(func() error {
			if raw.District.Equals(normalized.District) {
				return nil
			}
			return c.Set(ctx, raw.District.DocID(ctx).String(), normalized.District.String(), 0).Err()
		})
		return group.Wait()
	}
}
