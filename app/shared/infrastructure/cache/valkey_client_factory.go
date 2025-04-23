package cache

import (
	"transport-app/app/shared/configuration"

	"github.com/valkey-io/valkey-go"
)

func newValkeyClientFactory(conf configuration.Conf) (any, error) {
	opt, err := valkey.ParseURL(conf.CACHE_URL)
	if err != nil {
		return nil, err
	}
	c, err := valkey.NewClient(opt)
	return c, err
}
