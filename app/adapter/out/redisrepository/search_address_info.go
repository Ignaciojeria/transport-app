package redisrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/rediscache"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type SearchAddressInfo func(ctx context.Context, raw domain.AddressInfo) (domain.AddressInfo, error)

func init() {
	ioc.Registry(NewSearchAddressInfo, rediscache.NewRedisClient)
}

func NewSearchAddressInfo(c *redis.Client) SearchAddressInfo {
	return func(ctx context.Context, raw domain.AddressInfo) (domain.AddressInfo, error) {
		var result domain.AddressInfo

		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			key := raw.Province.DocID(ctx).String()
			val, err := c.Get(ctx, key).Result()
			if err == redis.Nil {
				return nil // no encontrada, pero no es error
			}
			if err != nil {
				return err
			}
			result.Province = domain.Province(val)
			return nil
		})

		group.Go(func() error {
			key := raw.State.DocID(ctx).String()
			val, err := c.Get(ctx, key).Result()
			if err == redis.Nil {
				return nil
			}
			if err != nil {
				return err
			}
			result.State = domain.State(val)
			return nil
		})

		group.Go(func() error {
			key := raw.District.DocID(ctx).String()
			val, err := c.Get(ctx, key).Result()
			if err == redis.Nil {
				return nil
			}
			if err != nil {
				return err
			}
			result.District = domain.District(val)
			return nil
		})

		err := group.Wait()
		return result, err
	}
}
