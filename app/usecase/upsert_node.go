package usecase

import (
	"context"
	"fmt"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertNode func(context.Context, domain.Origin) error

func init() {
	ioc.Registry(NewUpsertNode)
}

func NewUpsertNode() UpsertNode {
	return func(ctx context.Context, origin domain.Origin) error {
		fmt.Println("works")
		return nil
	}
}
