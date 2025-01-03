package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewSaveOrder,
		tidb.NewTIDBConnection,
		NewLoadOrderStatuses)
}

type SaveOrder func(
	context.Context,
	domain.Order) (domain.Order, error)

func NewSaveOrder(
	conn tidb.TIDBConnection,
	loadOrderSorderStatuses LoadOrderStatuses,
) SaveOrder {
	return func(ctx context.Context, to domain.Order) (domain.Order, error) {
		to.OrderStatus = loadOrderSorderStatuses().Available()
		//	table := mapper.MapOrderToTable(to)
		err := conn.Transaction(func(tx *gorm.DB) error {
			/*
				organizationID, err := ensureOrganizationExists(tx, table.Organization)
				if err != nil {
					return err
				}
				table.OrganizationID = organizationID*/
			/*
				commerceID, err := ensureCommerceExists(tx, organizationID, table.Commerce)
				if err != nil {
					return err
				}
				table.CommerceID = commerceID

				consumerID, err := ensureConsumerExists(tx, organizationID, table.Consumer)
				if err != nil {
					return err
				}
				table.ConsumerID = consumerID

				orderTypeID, err := ensureOrderTypeExists(tx, organizationID, table.OrderType)
				if err != nil {
					return err
				}
				table.OrderTypeID = orderTypeID
			*/
			/*
				if err := tx.Save(&table).Error; err != nil {
					return err
				}*/

			return nil
		})

		if err != nil {
			return domain.Order{}, err
		}

		return domain.Order{}, err
	}
}

func ensureOrganizationExists(tx *gorm.DB, organization table.Organization) (int64, error) {
	err := tx.
		Where("email = ?", organization.Email).
		//	Where("country = ?", organization.Country).
		First(&organization).Error
	if err == nil {
		return organization.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&organization).Error; err != nil {
			return 0, err
		}
		return organization.ID, nil
	}
	return 0, err
}

func ensureCommerceExists(tx *gorm.DB, organizationID int64, commerce table.Commerce) (int, error) {
	var com table.Commerce
	err := tx.Where("name = ? AND organization_id = ?", commerce.Name, organizationID).First(&com).Error
	if err == nil {
		return com.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		com = table.Commerce{Name: commerce.Name, OrganizationID: organizationID}
		if err := tx.Create(&com).Error; err != nil {
			return 0, err
		}
		return com.ID, nil
	}
	return 0, err
}

func ensureConsumerExists(tx *gorm.DB, organizationID int64, consumer table.Consumer) (int, error) {
	var con table.Consumer
	err := tx.Where("name = ? AND organization_id = ?", consumer.Name, organizationID).First(&con).Error
	if err == nil {
		return con.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		con = table.Consumer{Name: consumer.Name, OrganizationID: organizationID}
		if err := tx.Create(&con).Error; err != nil {
			return 0, err
		}
		return con.ID, nil
	}
	return 0, err
}

func ensureOrderTypeExists(tx *gorm.DB, organizationID int64, orderType table.OrderType) (int64, error) {
	var ot table.OrderType
	err := tx.Where("type = ? AND organization_id = ?", orderType.Type, organizationID).First(&ot).Error
	if err == nil {
		return ot.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ot = table.OrderType{Type: orderType.Type, Description: orderType.Description, OrganizationID: organizationID}
		if err := tx.Create(&ot).Error; err != nil {
			return 0, err
		}
		return ot.ID, nil
	}
	return 0, err
}
