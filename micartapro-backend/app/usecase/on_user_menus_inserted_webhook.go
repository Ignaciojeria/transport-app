package usecase

import (
	"context"
	"fmt"
	"micartapro/app/events"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OnUserMenusInsertedWebhook func(ctx context.Context, input events.UserMenusInsertedWebhook) error

func init() {
	ioc.Registry(NewOnUserMenusInsertedWebhook)
}

func NewOnUserMenusInsertedWebhook() OnUserMenusInsertedWebhook {
	return func(ctx context.Context, events events.UserMenusInsertedWebhook) error {
		fmt.Println("TODO : BUILD menuInteractionRequest with default values & publish it")
		return nil
	}
}
