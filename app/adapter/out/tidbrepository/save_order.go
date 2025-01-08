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

		type QueryResult struct {
			OrganizationCountryID int64
			CommerceID            int64
			ConsumerID            int64
			OrderTypeID           int64
			OriginContactID       int64
			DestinationContactID  int64
			OriginAddressID       int64
			DestinationAddressID  int64
			OriginNodeInfoID      int64
			DestinationNodeInfoID int64
		}

		var result QueryResult

		err := conn.Raw(`
WITH 
  org_country AS (
    SELECT 
      organization_countries.id AS organization_country_id
    FROM 
      organization_countries
    INNER JOIN api_keys ON api_keys.organization_id = organization_countries.organization_id
    WHERE 
      organization_countries.country = ? AND api_keys.key = ?
  )
SELECT 
  org_country.organization_country_id AS organization_country_id,
  c.id AS commerce_id,
  con.id AS consumer_id,
  ot.id AS order_type_id,
  o_ct.id AS origin_contact_id,
  d_ct.id AS destination_contact_id,
  o_ai.id AS origin_address_id,
  d_ai.id AS destination_address_id,
  o_ni.id AS origin_node_info_id,
  d_ni.id AS destination_node_info_id
FROM 
  org_country
LEFT JOIN commerces c ON c.organization_country_id = org_country.organization_country_id AND c.name = ?
LEFT JOIN consumers con ON con.organization_country_id = org_country.organization_country_id AND con.name = ?
LEFT JOIN order_types ot ON ot.organization_country_id = org_country.organization_country_id AND ot.type = ?
LEFT JOIN contacts o_ct ON o_ct.organization_country_id = org_country.organization_country_id AND o_ct.full_name = ? AND o_ct.email = ? AND o_ct.phone = ?
LEFT JOIN contacts d_ct ON d_ct.organization_country_id = org_country.organization_country_id AND d_ct.full_name = ? AND d_ct.email = ? AND d_ct.phone = ?
LEFT JOIN address_infos o_ai ON o_ai.organization_country_id = org_country.organization_country_id AND o_ai.raw_address = ?
LEFT JOIN address_infos d_ai ON d_ai.organization_country_id = org_country.organization_country_id AND d_ai.raw_address = ?
LEFT JOIN node_infos o_ni ON o_ni.organization_country_id = org_country.organization_country_id AND o_ni.reference_id = ?
LEFT JOIN node_infos d_ni ON d_ni.organization_country_id = org_country.organization_country_id AND d_ni.reference_id = ?;

		`,
			to.Organization.Country.Alpha2(),
			to.Organization.Key,
			to.BusinessIdentifiers.Commerce,
			to.BusinessIdentifiers.Consumer,
			to.OrderType.Type,
			to.Origin.AddressInfo.Contact.FullName,
			to.Origin.AddressInfo.Contact.Email,
			to.Origin.AddressInfo.Contact.Phone,
			to.Destination.AddressInfo.Contact.FullName,
			to.Destination.AddressInfo.Contact.Email,
			to.Destination.AddressInfo.Contact.Phone,
			to.Origin.AddressInfo.RawAddress(),
			to.Destination.AddressInfo.RawAddress(),
			to.Origin.NodeInfo.ReferenceID,
			to.Destination.NodeInfo.ReferenceID,
		).Scan(&result).Error

		if err != nil {
			return domain.Order{}, err
		}

		return domain.Order{}, conn.Transaction(func(tx *gorm.DB) error {
			orderTable := mapper.MapOrderToTable(to)
			if result.CommerceID == 0 {
				orderTable.Commerce.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.Commerce).Error; err != nil {
					return err
				}
			}
			if result.ConsumerID == 0 {
				orderTable.Consumer.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.Consumer).Error; err != nil {
					return err
				}
			}

			if result.OriginContactID == 0 {
				orderTable.OriginContact.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.OriginContact).Error; err != nil {
					return err
				}
			}

			if result.DestinationContactID == 0 && !to.IsOriginAndDestinationContactEqual() {
				orderTable.DestinationContact.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.DestinationContact).Error; err != nil {
					return err
				}
			}

			if result.OriginAddressID == 0 {
				orderTable.OriginAddressInfo.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.OriginAddressInfo).Error; err != nil {
					return err
				}
			}

			if result.DestinationAddressID == 0 && !to.IsOriginAndDestinationAddressEqual() {
				orderTable.DestinationAddressInfo.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.DestinationAddressInfo).Error; err != nil {
					return err
				}
			}

			/*
				if result.OriginNodeInfoID == 0 {
					orderTable.OriginNodeInfo.OrganizationCountryID = result.OrganizationCountryID
					orderTable.Origin
					if err := tx.Save(&orderTable.OriginNodeInfo).Error; err != nil {
						return err
					}
				}
				if result.DestinationNodeInfoID == 0 {
					orderTable.DestinationNodeInfo.OrganizationCountryID = result.OrganizationCountryID
					if err := tx.Save(&orderTable.DestinationNodeInfo).Error; err != nil {
						return err
					}
				}
			*/
			return nil
		})
	}
}

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

func ensureCommerceExists(tx *gorm.DB, organizationCountryID int64, commerce table.Commerce) (int64, error) {
	var com table.Commerce
	err := tx.Where("name = ? AND organization_country_id = ?", commerce.Name, organizationCountryID).First(&com).Error
	if err == nil {
		return com.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		com = table.Commerce{Name: commerce.Name, OrganizationCountryID: organizationCountryID}
		if err := tx.Create(&com).Error; err != nil {
			return 0, err
		}
		return com.ID, nil
	}
	return 0, err
}

func ensureConsumerExists(tx *gorm.DB, organizationCountryID int64, consumer table.Consumer) (int64, error) {
	var con table.Consumer
	err := tx.Where("name = ? AND organization_country_id = ?", consumer.Name, organizationCountryID).First(&con).Error
	if err == nil {
		return con.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		con = table.Consumer{Name: consumer.Name, OrganizationCountryID: organizationCountryID}
		if err := tx.Create(&con).Error; err != nil {
			return 0, err
		}
		return con.ID, nil
	}
	return 0, err
}

func ensureOrderTypeExists(tx *gorm.DB, organizationCountryID int64, orderType table.OrderType) (int64, error) {
	var ot table.OrderType
	err := tx.Where("type = ? AND organization_country_id = ?", orderType.Type, organizationCountryID).First(&ot).Error
	if err == nil {
		return ot.ID, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ot = table.OrderType{Type: orderType.Type, Description: orderType.Description, OrganizationCountryID: organizationCountryID}
		if err := tx.Create(&ot).Error; err != nil {
			return 0, err
		}
		return ot.ID, nil
	}
	return 0, err
}
