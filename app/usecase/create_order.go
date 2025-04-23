package usecase

import (
	"context"
	"transport-app/app/adapter/out/geocoding"
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
		tidbrepository.NewUpsertOrderReferences,
		tidbrepository.NewUpsertOrderPackages,
		geocoding.NewGeocodingStrategy,
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
	upsertOrderReferences tidbrepository.UpsertOrderReferences,
	upsertOrderPackages tidbrepository.UpsertOrderPackages,
	geocode geocoding.GeocodingStrategy,
) CreateOrder {
	return func(ctx context.Context, inOrder domain.Order) error {
		inOrder.OrderStatus = loadOrderStatuses().Available()

		normalizationGroup, ctx := errgroup.WithContext(ctx)
		normalizationGroup.Go(func() error {
			return inOrder.Origin.AddressInfo.NormalizeAndGeocode(
				ctx,
				geocode,
			)
		})

		normalizationGroup.Go(func() error {
			return inOrder.Destination.AddressInfo.NormalizeAndGeocode(
				ctx,
				geocode,
			)
		})

		if err := normalizationGroup.Wait(); err != nil {
			return err
		}

		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			return upsertOrderHeaders(ctx, inOrder.Headers)
		})

		group.Go(func() error {
			return upsertContact(ctx, inOrder.Origin.AddressInfo.Contact)
		})

		group.Go(func() error {
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
			return upsertPackages(ctx, inOrder.Packages,
				inOrder.
					ReferenceID.
					String())
		})

		group.Go(func() error {
			return upsertOrderReferences(ctx, inOrder)
		})

		group.Go(func() error {
			return upsertOrderPackages(ctx, inOrder)
		})

		group.Go(func() error {
			return upsertOrder(ctx, inOrder)
		})

		return group.Wait()
	}
}
