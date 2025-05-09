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
		tidbrepository.NewUpsertDeliveryUnits,
		tidbrepository.NewUpsertOrderType,
		tidbrepository.NewUpsertOrder,
		tidbrepository.NewUpsertOrderReferences,
		tidbrepository.NewUpsertOrderDeliveryUnits,
		geocoding.NewGeocodingStrategy,
	)
}

func NewCreateOrder(
	upsertOrderHeaders tidbrepository.UpsertOrderHeaders,
	loadOrderStatuses tidbrepository.LoadOrderStatuses,
	upsertContact tidbrepository.UpsertContact,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
	upsertDeliveryUnits tidbrepository.UpsertDeliveryUnits,
	upsertOrderType tidbrepository.UpsertOrderType,
	upsertOrder tidbrepository.UpsertOrder,
	upsertOrderReferences tidbrepository.UpsertOrderReferences,
	upsertOrderDeliveryUnits tidbrepository.UpsertOrderDeliveryUnits,
	geocode geocoding.GeocodingStrategy,
) CreateOrder {
	return func(ctx context.Context, inOrder domain.Order) error {
		inOrder.OrderStatus = loadOrderStatuses().Available()

		normalizationGroup, group1Ctx := errgroup.WithContext(ctx)

		normalizationGroup.Go(func() error {
			return inOrder.Origin.AddressInfo.NormalizeAndGeocode(
				group1Ctx,
				geocode,
			)
		})

		normalizationGroup.Go(func() error {
			return inOrder.Destination.AddressInfo.NormalizeAndGeocode(
				group1Ctx,
				geocode,
			)
		})

		if err := normalizationGroup.Wait(); err != nil {
			return err
		}

		group, group2Ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			return upsertOrderHeaders(group2Ctx, inOrder.Headers)
		})

		group.Go(func() error {
			return upsertContact(group2Ctx, inOrder.Origin.AddressInfo.Contact)
		})

		group.Go(func() error {
			return upsertContact(group2Ctx, inOrder.Destination.AddressInfo.Contact)
		})

		group.Go(func() error {
			return upsertAddressInfo(group2Ctx, inOrder.Origin.AddressInfo)
		})

		group.Go(func() error {
			return upsertAddressInfo(group2Ctx, inOrder.Destination.AddressInfo)
		})

		group.Go(func() error {
			return upsertNodeInfo(group2Ctx, inOrder.Origin)
		})

		group.Go(func() error {
			return upsertNodeInfo(group2Ctx, inOrder.Destination)
		})

		group.Go(func() error {
			return upsertOrderType(group2Ctx, inOrder.OrderType)
		})

		group.Go(func() error {
			return upsertDeliveryUnits(group2Ctx, inOrder.Packages,
				inOrder.
					ReferenceID.
					String())
		})

		group.Go(func() error {
			return upsertOrderReferences(group2Ctx, inOrder)
		})

		group.Go(func() error {
			return upsertOrderDeliveryUnits(group2Ctx, inOrder)
		})

		group.Go(func() error {
			return upsertOrder(group2Ctx, inOrder)
		})

		return group.Wait()
	}
}
