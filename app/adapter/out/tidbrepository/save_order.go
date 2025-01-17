package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewSaveOrder,
		tidb.NewTIDBConnection,
		NewLoadOrderStatuses,
	)
}

type SaveOrder func(
	ctx context.Context, orderToCreate domain.Order) (domain.Order, error)

func NewSaveOrder(
	conn tidb.TIDBConnection,
	loadOrderStatuses LoadOrderStatuses,
) SaveOrder {
	return func(ctx context.Context, orderToCreate domain.Order) (domain.Order, error) {
		orderTable := mapper.MapOrderToTable(orderToCreate)

		err := conn.Transaction(func(tx *gorm.DB) error {
			// Evaluar si el origen y el destino son iguales para AddressInfo
			if orderToCreate.IsOriginAndDestinationAddressEqual() && orderTable.OriginAddressInfo.ID == 0 {
				// Guardar el origen
				err := tx.Save(&orderTable.OriginAddressInfo).Error
				if err != nil {
					return err
				}
				// Setear el origen al destino
				orderTable.DestinationAddressInfo.ID = orderTable.OriginAddressInfo.ID
			}

			// Evaluar si el origen y el destino son iguales para Contact
			if orderToCreate.IsOriginAndDestinationContactEqual() && orderTable.OriginContact.ID == 0 {
				// Guardar el contacto de origen
				err := tx.Save(&orderTable.OriginContact).Error
				if err != nil {
					return err
				}

				// Setear el contacto de origen al destino
				orderTable.DestinationContact.ID = orderTable.OriginContact.ID
			}

			// Evaluar si el ReferenceID del nodo de origen y destino son iguales
			if orderToCreate.IsOriginAndDestinationNodeEqual() && orderTable.OriginNodeInfo.ID == 0 {
				// Guardar el nodo de origen
				err := tx.Save(&orderTable.OriginNodeInfo).Error
				if err != nil {
					return err
				}

				// Setear el nodo de origen al destino
				orderTable.DestinationNodeInfo.ID = orderTable.OriginNodeInfo.ID
			}

			// Guardar la orden completa
			return tx.Save(&orderTable).Error
		})

		if err != nil {
			return domain.Order{}, err
		}

		return domain.Order{}, nil
	}
}
