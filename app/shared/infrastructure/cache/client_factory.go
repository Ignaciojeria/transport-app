package cache

import (
	"errors"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewCacheClientFactory,
		configuration.NewConf)
}

func NewCacheClientFactory(conf configuration.Conf) (any, error) {
	if conf.CACHE_STRATEGY == "valkey" {
		return newValkeyClientFactory(conf)
	}
	if conf.CACHE_STRATEGY == "redis" {
		return newRedisClientFactory(conf)
	}
	return nil, errors.New("unimplemented cache")
}
