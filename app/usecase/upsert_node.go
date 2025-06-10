package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertNode func(ctx context.Context, input domain.NodeInfo) error

func init() {
	ioc.Registry(
		NewUpsertNode,
		tidbrepository.NewUpsertNodeInfo,
		tidbrepository.NewUpsertContact,
		tidbrepository.NewUpsertAddressInfo,
		tidbrepository.NewUpsertNodeType,
		tidbrepository.NewUpsertNodeInfoHeaders,
	)
}

func NewUpsertNode(
	upsertNodeInfo tidbrepository.UpsertNodeInfo,
	upsertContact tidbrepository.UpsertContact,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	upsertNodeType tidbrepository.UpsertNodeType,
	upsertNodeInfoHeaders tidbrepository.UpsertNodeInfoHeaders,
) UpsertNode {
	return func(ctx context.Context, input domain.NodeInfo) error {

		// Actualizar o insertar el tipo de nodo
		if err := upsertNodeType(ctx, input.NodeType); err != nil {
			return err
		}

		// Actualizar o insertar el contacto
		if err := upsertContact(ctx, input.AddressInfo.Contact); err != nil {
			return err
		}

		// Actualizar o insertar la información de dirección
		if err := upsertAddressInfo(ctx, input.AddressInfo); err != nil {
			return err
		}

		// Finalmente, actualizar o insertar el nodo
		return upsertNodeInfo(ctx, input)
	}
}
