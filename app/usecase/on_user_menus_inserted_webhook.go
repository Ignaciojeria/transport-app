package usecase

import (
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnUserMenusInsertedWebhook func(ctx context.Context, input interface{}) (interface{}, error)

func init() {
	ioc.Registry(NewOnUserMenusInsertedWebhook)
}

func NewOnUserMenusInsertedWebhook() OnUserMenusInsertedWebhook {
	return func(ctx context.Context, input interface{}) (interface{}, error) {
		return input, nil
	}
}
