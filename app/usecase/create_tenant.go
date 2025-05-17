package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"golang.org/x/sync/errgroup"
)

type CreateTenant func(ctx context.Context, org domain.Tenant) error

func init() {
	ioc.Registry(
		NewCreateTenant,
		tidbrepository.NewUpsertOrderHeaders,
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertAddressInfo,
		tidbrepository.NewUpsertNodeInfo,
		tidbrepository.NewUpsertDeliveryUnits,
		tidbrepository.NewUpsertOrderType,
		tidbrepository.NewUpsertOrder,
		tidbrepository.NewUpsertOrderReferences,
		tidbrepository.NewUpsertOrderDeliveryUnits,
		tidbrepository.NewSaveTenant,
		tidbrepository.NewUpsertVehicleCategory,
		tidbrepository.NewUpsertCarrier,
		tidbrepository.NewUpsertDriver,
		tidbrepository.NewUpsertVehicle,
		tidbrepository.NewUpsertVehicleHeaders,
		tidbrepository.NewUpsertPlan,
		tidbrepository.NewUpsertPlanHeaders,
		tidbrepository.NewUpsertRoute,
		tidbrepository.NewUpsertState,
		tidbrepository.NewUpsertProvince,
		tidbrepository.NewUpsertDistrict,
	)
}

func NewCreateTenant(
	upsertOrderHeaders tidbrepository.UpsertOrderHeaders,
	upsertContact tidbrepository.UpsertContact,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
	upsertDeliveryUnits tidbrepository.UpsertDeliveryUnits,
	upsertOrderType tidbrepository.UpsertOrderType,
	upsertOrder tidbrepository.UpsertOrder,
	upsertOrderReferences tidbrepository.UpsertOrderReferences,
	upsertOrderDeliveryUnits tidbrepository.UpsertOrderDeliveryUnits,
	saveTenant tidbrepository.SaveTenant,
	upsertVehicleCategory tidbrepository.UpsertVehicleCategory,
	upsertCarrier tidbrepository.UpsertCarrier,
	upsertDriver tidbrepository.UpsertDriver,
	upsertVehicle tidbrepository.UpsertVehicle,
	upsertVehicleHeaders tidbrepository.UpsertVehicleHeaders,
	upsertPlan tidbrepository.UpsertPlan,
	upsertPlanHeaders tidbrepository.UpsertPlanHeaders,
	upsertRoute tidbrepository.UpsertRoute,
	upsertState tidbrepository.UpsertState,
	upsertProvince tidbrepository.UpsertProvince,
	upsertDistrict tidbrepository.UpsertDistrict,
) CreateTenant {
	return func(ctx context.Context, org domain.Tenant) error {
		_, err := saveTenant(ctx, org)
		if err != nil {
			return err
		}
		group, groupCtx := errgroup.WithContext(ctx)

		group.Go(func() error {
			return upsertOrderHeaders(groupCtx, domain.Headers{})
		})

		group.Go(func() error {
			return upsertContact(groupCtx, domain.Contact{})
		})

		group.Go(func() error {
			return upsertAddressInfo(groupCtx, domain.AddressInfo{})
		})

		group.Go(func() error {
			return upsertNodeInfo(groupCtx, domain.NodeInfo{})
		})

		group.Go(func() error {
			return upsertDeliveryUnits(groupCtx, []domain.Package{}, "")
		})

		group.Go(func() error {
			return upsertOrderType(groupCtx, domain.OrderType{})
		})

		group.Go(func() error {
			return upsertOrderReferences(groupCtx, domain.Order{})
		})

		group.Go(func() error {
			return upsertOrderDeliveryUnits(groupCtx, domain.Order{})
		})

		group.Go(func() error {
			return upsertVehicleCategory(groupCtx, domain.VehicleCategory{})
		})

		group.Go(func() error {
			return upsertCarrier(groupCtx, domain.Carrier{})
		})

		group.Go(func() error {
			return upsertDriver(groupCtx, domain.Driver{})
		})

		group.Go(func() error {
			return upsertVehicle(groupCtx, domain.Vehicle{})
		})

		group.Go(func() error {
			return upsertVehicleHeaders(groupCtx, domain.Headers{})
		})

		group.Go(func() error {
			return upsertPlan(groupCtx, domain.Plan{})
		})

		group.Go(func() error {
			return upsertPlanHeaders(groupCtx, domain.Headers{})
		})
		group.Go(func() error {
			return upsertRoute(groupCtx, domain.Route{}, "")
		})

		group.Go(func() error {
			return upsertState(groupCtx, domain.State(""))
		})

		group.Go(func() error {
			return upsertProvince(groupCtx, domain.Province(""))
		})

		group.Go(func() error {
			return upsertDistrict(groupCtx, domain.District(""))
		})

		return group.Wait()
	}
}
