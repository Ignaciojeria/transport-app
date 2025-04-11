package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"golang.org/x/sync/errgroup"
)

type CreateOrder func(ctx context.Context, input domain.Order) error

func init() {
	ioc.Registry(
		NewCreateOrder,
		tidbrepository.NewUpsertOrderHeaders,
		tidbrepository.NewLoadOrderStatuses,
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertAddressInfo,
		tidbrepository.NewUpsertNodeInfo,
		tidbrepository.NewUpsertPackages,
		tidbrepository.NewUpsertOrderType,
		tidbrepository.NewUpsertOrder,
	)
}

func NewCreateOrder(
	upsertOrderHeaders tidbrepository.UpsertOrderHeaders,
	loadOrderStatuses tidbrepository.LoadOrderStatuses,
	upsertContact tidbrepository.UpsertContact,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
	upsertPackages tidbrepository.UpsertPackages,
	upsertOrderType tidbrepository.UpsertOrderType,
	upsertOrder tidbrepository.UpsertOrder,
) CreateOrder {
	return func(ctx context.Context, inOrder domain.Order) error {
		inOrder.OrderStatus = loadOrderStatuses().Available()

		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			return upsertOrderHeaders(ctx, inOrder.Headers)
		})

		group.Go(func() error {
			inOrder.Origin.AddressInfo.Normalize()
			return upsertContact(ctx, inOrder.Origin.AddressInfo.Contact)
		})

		group.Go(func() error {
			inOrder.Destination.AddressInfo.Normalize()
			return upsertContact(ctx, inOrder.Destination.AddressInfo.Contact)
		})

		group.Go(func() error {
			return upsertAddressInfo(ctx, inOrder.Origin.AddressInfo)
		})

		group.Go(func() error {
			return upsertAddressInfo(ctx, inOrder.Destination.AddressInfo)
		})

		group.Go(func() error {
			return upsertNodeInfo(ctx, inOrder.Origin)
		})

		group.Go(func() error {
			return upsertNodeInfo(ctx, inOrder.Destination)
		})

		group.Go(func() error {
			return upsertOrderType(ctx, inOrder.OrderType)
		})

		group.Go(func() error {
			return upsertPackages(ctx, inOrder.Packages)
		})

		group.Go(func() error {
			return upsertOrder(ctx, inOrder)
		})

		return group.Wait()
	}
}
