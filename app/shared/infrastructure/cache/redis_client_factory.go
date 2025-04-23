package cache

import (
	"transport-app/app/shared/configuration"

	"github.com/redis/go-redis/v9"
)

func newRedisClientFactory(conf configuration.Conf) (any, error) {
	opt, err := redis.ParseURL(conf.CACHE_URL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return client, nil
}
