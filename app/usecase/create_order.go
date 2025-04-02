package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateOrder func(ctx context.Context, input domain.Order) (domain.Order, error)

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
	return func(ctx context.Context, inOrder domain.Order) (domain.Order, error) {
		inOrder.OrderStatus = loadOrderStatuses().Available()

		inOrder.Headers.Organization = inOrder.Organization
		err := upsertOrderHeaders(ctx, inOrder.Headers)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Origin.AddressInfo.Contact.Organization = inOrder.Organization
		err = upsertContact(ctx, inOrder.Origin.AddressInfo.Contact)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Destination.AddressInfo.Contact.Organization = inOrder.Organization
		err = upsertContact(ctx, inOrder.Destination.AddressInfo.Contact)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Origin.AddressInfo.Organization = inOrder.Organization
		err = upsertAddressInfo(ctx, inOrder.Origin.AddressInfo)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Destination.AddressInfo.Organization = inOrder.Organization
		err = upsertAddressInfo(ctx, inOrder.Destination.AddressInfo)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Origin.Organization = inOrder.Organization
		originNodeInfo, err := upsertNodeInfo(ctx, inOrder.Origin)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Destination.Organization = inOrder.Organization
		destinationNodeInfo, err := upsertNodeInfo(ctx, inOrder.Destination)
		if err != nil {
			return domain.Order{}, err
		}
		inOrder.OrderType.Organization = inOrder.Organization
		err = upsertOrderType(ctx, inOrder.OrderType)
		if err != nil {
			return domain.Order{}, err
		}

		pcks, err := upsertPackages(ctx, inOrder.Packages, inOrder.Organization)
		if err != nil {
			return domain.Order{}, err
		}
		//inOrder.Headers = orderHeaders
		//inOrder.Origin.AddressInfo.Contact = originContact
		//inOrder.Destination.AddressInfo.Contact = destinationContact
		//inOrder.OrderType = orderType
		inOrder.Origin = originNodeInfo
		inOrder.Destination = destinationNodeInfo
		//inOrder.Origin.AddressInfo = originAddressInfo
		//inOrder.Destination.AddressInfo = destinationAddressInfo

		inOrder.Packages = pcks
		return upsertOrder(ctx, inOrder)
	}
}
