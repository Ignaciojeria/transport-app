package rediscache

import (
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/redis/go-redis/v9"
)

func init() {
	ioc.Registry(NewRedisClient, configuration.NewConf)
}
func NewRedisClient(conf configuration.Conf) (*redis.Client, error) {
	opt, err := redis.ParseURL(conf.REDIS_URL)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
}
