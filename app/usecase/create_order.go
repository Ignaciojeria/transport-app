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
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertAddressInfo,
		tidbrepository.NewUpsertDeliveryUnits,
		tidbrepository.NewUpsertOrderType,
		tidbrepository.NewUpsertOrder,
		tidbrepository.NewUpsertOrderReferences,
		tidbrepository.NewUpsertOrderDeliveryUnits,
		tidbrepository.NewUpsertDeliveryUnitsHistory,
		tidbrepository.NewUpsertSizeCategory,
		tidbrepository.NewUpsertDeliveryUnitsLabels,
		tidbrepository.NewUpsertSkill,
		tidbrepository.NewUpsertDeliveryUnitsSkills,
		geocoding.NewGeocodingStrategy,
	)
}

func NewCreateOrder(
	upsertOrderHeaders tidbrepository.UpsertOrderHeaders,
	upsertContact tidbrepository.UpsertContact,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertDeliveryUnits tidbrepository.UpsertDeliveryUnits,
	upsertOrderType tidbrepository.UpsertOrderType,
	upsertOrder tidbrepository.UpsertOrder,
	upsertOrderReferences tidbrepository.UpsertOrderReferences,
	upsertOrderDeliveryUnits tidbrepository.UpsertOrderDeliveryUnits,
	upsertDeliveryUnitsHistory tidbrepository.UpsertDeliveryUnitsHistory,
	upsertSizeCategory tidbrepository.UpsertSizeCategory,
	upsertDeliveryUnitsLabels tidbrepository.UpsertDeliveryUnitsLabels,
	upsertSkill tidbrepository.UpsertSkill,
	upsertDeliveryUnitsSkills tidbrepository.UpsertDeliveryUnitsSkills,
	geocode geocoding.GeocodingStrategy,
) CreateOrder {
	return func(ctx context.Context, inOrder domain.Order) error {
		normalizationGroup, group1Ctx := errgroup.WithContext(ctx)
		inOrder.Origin.AddressInfo.ToLowerAndRemovePunctuation()
		inOrder.Destination.AddressInfo.ToLowerAndRemovePunctuation()
		inOrder.AssignIndexesIfNoLPN()
		normalizationGroup.Go(func() error {
			if inOrder.Origin.AddressInfo.Equals(group1Ctx, inOrder.Destination.AddressInfo) {
				return nil
			}
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
			if inOrder.Origin.AddressInfo.Contact.Equals(group2Ctx, inOrder.Destination.AddressInfo.Contact) {
				return nil
			}
			return upsertContact(group2Ctx, inOrder.Origin.AddressInfo.Contact)
		})

		group.Go(func() error {
			return upsertContact(group2Ctx, inOrder.Destination.AddressInfo.Contact)
		})

		group.Go(func() error {
			if inOrder.Origin.AddressInfo.Equals(group1Ctx, inOrder.Destination.AddressInfo) {
				return nil
			}
			return upsertAddressInfo(group2Ctx, inOrder.Origin.AddressInfo)
		})

		group.Go(func() error {
			return upsertAddressInfo(group2Ctx, inOrder.Destination.AddressInfo)
		})

		group.Go(func() error {
			return upsertOrderType(group2Ctx, inOrder.OrderType)
		})

		group.Go(func() error {
			return upsertDeliveryUnits(group2Ctx, inOrder.DeliveryUnits)
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

		group.Go(func() error {
			return upsertDeliveryUnitsLabels(group2Ctx, inOrder)
		})

		group.Go(func() error {
			return upsertDeliveryUnitsSkills(group2Ctx, inOrder)
		})

		group.Go(func() error {
			processedSkills := make(map[domain.Skill]struct{})
			for _, deliveryUnit := range inOrder.DeliveryUnits {
				for _, skill := range deliveryUnit.Skills {
					if _, exists := processedSkills[skill]; !exists {
						processedSkills[skill] = struct{}{}
						if err := upsertSkill(group2Ctx, skill); err != nil {
							return err
						}
					}
				}
			}
			return nil
		})

		group.Go(func() error {
			plan := domain.Plan{
				Routes: []domain.Route{
					{
						Orders: []domain.Order{inOrder},
					},
				},
			}
			return upsertDeliveryUnitsHistory(group2Ctx, plan)
		})

		// Upsert size categories for each delivery unit
		for _, deliveryUnit := range inOrder.DeliveryUnits {
			du := deliveryUnit // Create a new variable to avoid closure issues
			group.Go(func() error {
				return upsertSizeCategory(group2Ctx, du.SizeCategory)
			})
		}

		return group.Wait()
	}
}
