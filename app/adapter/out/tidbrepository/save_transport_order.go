package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/mapper"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(NewSaveTransportOrder, tidb.NewTIDBConnection)
}

type SaveTransportOrder func(context.Context, domain.TransportOrder) (domain.TransportOrder, error)

func NewSaveTransportOrder(conn tidb.TIDBConnection) SaveTransportOrder {
	return func(ctx context.Context, to domain.TransportOrder) (domain.TransportOrder, error) {
		table := mapper.MapTransportOrderToTable(to)

		err := conn.Transaction(func(tx *gorm.DB) error {
			organizationID, err := ensureOrganizationExists(tx, table.Organization)
			if err != nil {
				return err
			}
			table.OrganizationID = organizationID

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

			/*
				if err := tx.Save(&table).Error; err != nil {
					return err
				}*/

			return nil
		})

		if err != nil {
			return domain.TransportOrder{}, err
		}

		return domain.TransportOrder{}, err
	}
}

func ensureOrganizationExists(tx *gorm.DB, organization table.Organization) (int64, error) {
	var org table.Organization
	err := tx.Where("name = ?", organization.Name).First(&org).Error
	if err == nil {
		return org.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		org = table.Organization{Name: organization.Name}
		if err := tx.Create(&org).Error; err != nil {
			return 0, err
		}
		return org.ID, nil
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
