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
		tidbrepository.NewLoadOrderStatuses,
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertAddressInfo,
		tidbrepository.NewUpsertNodeInfo,
	)
}

func NewCreateOrder(
	loadOrderStatuses tidbrepository.LoadOrderStatuses,
	upsertContact tidbrepository.UpsertContact,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
) CreateOrder {
	return func(ctx context.Context, inOrder domain.Order) (domain.Order, error) {
		inOrder.OrderStatus = loadOrderStatuses().Available()

		inOrder.Origin.AddressInfo.Contact.Organization = inOrder.Organization
		originContact, err := upsertContact(ctx, inOrder.Origin.AddressInfo.Contact)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Destination.AddressInfo.Contact.Organization = inOrder.Organization
		destinationContact, err := upsertContact(ctx, inOrder.Destination.AddressInfo.Contact)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Origin.AddressInfo.Organization = inOrder.Organization
		originAddressInfo, err := upsertAddressInfo(ctx, inOrder.Origin.AddressInfo)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Destination.AddressInfo.Organization = inOrder.Organization
		destinationAddressInfo, err := upsertAddressInfo(ctx, inOrder.Destination.AddressInfo)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Origin.Organization = inOrder.Organization
		originNodeInfo, err := upsertNodeInfo(ctx, inOrder.Origin)
		if err != nil {
			return domain.Order{}, err
		}

		inOrder.Origin = originNodeInfo
		inOrder.Origin.AddressInfo = originAddressInfo
		inOrder.Origin.AddressInfo.Contact = originContact
		inOrder.Destination.AddressInfo = destinationAddressInfo
		inOrder.Destination.AddressInfo.Contact = destinationContact
		return inOrder, nil
	}
}
