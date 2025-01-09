package tidbrepository

import (
	"context"
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

		to.OrderStatus = loadOrderSorderStatuses().Available()

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
		orderTable := mapper.MapOrderToTable(to)
		orderTable.CommerceID = result.CommerceID
		orderTable.ConsumerID = result.ConsumerID
		orderTable.OriginContactID = result.OriginContactID
		orderTable.DestinationContactID = result.DestinationContactID
		orderTable.OriginAddressInfoID = result.OriginAddressID
		orderTable.DestinationAddressInfoID = result.DestinationAddressID
		orderTable.OriginNodeInfoID = result.OriginNodeInfoID
		orderTable.DestinationNodeInfoID = result.DestinationNodeInfoID
		orderTable.OrderTypeID = result.OrderTypeID

		return domain.Order{}, conn.Transaction(func(tx *gorm.DB) error {
			// Guardar entidades que no existen y actualizar relaciones en orderTable
			if result.CommerceID == 0 {
				orderTable.Commerce.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.Commerce).Error; err != nil {
					return err
				}
				orderTable.CommerceID = orderTable.Commerce.ID
			}

			if result.ConsumerID == 0 {
				orderTable.Consumer.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.Consumer).Error; err != nil {
					return err
				}
				orderTable.ConsumerID = orderTable.Consumer.ID
			}

			if result.OriginContactID == 0 {
				orderTable.OriginContact.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.OriginContact).Error; err != nil {
					return err
				}
				orderTable.OriginContactID = orderTable.OriginContact.ID
			}

			if result.DestinationContactID == 0 {
				if to.IsOriginAndDestinationContactEqual() {
					orderTable.DestinationContactID = orderTable.OriginContactID
				} else {
					orderTable.DestinationContact.OrganizationCountryID = result.OrganizationCountryID
					if err := tx.Save(&orderTable.DestinationContact).Error; err != nil {
						return err
					}
					orderTable.DestinationContactID = orderTable.DestinationContact.ID
				}
			}

			if result.OriginAddressID == 0 {
				orderTable.OriginAddressInfo.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Save(&orderTable.OriginAddressInfo).Error; err != nil {
					return err
				}
				orderTable.OriginAddressInfoID = orderTable.OriginAddressInfo.ID
			}

			if result.DestinationAddressID == 0 {
				if to.IsOriginAndDestinationAddressEqual() {
					orderTable.DestinationAddressInfoID = orderTable.OriginAddressInfoID
				} else {
					orderTable.DestinationAddressInfo.OrganizationCountryID = result.OrganizationCountryID
					if err := tx.Save(&orderTable.DestinationAddressInfo).Error; err != nil {
						return err
					}
					orderTable.DestinationAddressInfoID = orderTable.DestinationAddressInfo.ID
				}
			}

			if result.OriginNodeInfoID == 0 {
				orderTable.OriginNodeInfo.OrganizationCountryID = result.OrganizationCountryID
				orderTable.OriginNodeInfo.AddressID = orderTable.OriginAddressInfoID
				if err := tx.Save(&orderTable.OriginNodeInfo).Error; err != nil {
					return err
				}
				orderTable.OriginNodeInfoID = orderTable.OriginNodeInfo.ID
			}

			if result.DestinationNodeInfoID == 0 {
				if to.IsOriginAndDestinationNodeEqual() {
					orderTable.DestinationNodeInfoID = orderTable.OriginNodeInfoID
				} else {
					orderTable.DestinationNodeInfo.OrganizationCountryID = result.OrganizationCountryID
					orderTable.DestinationNodeInfo.AddressID = orderTable.DestinationAddressInfoID
					if err := tx.Save(&orderTable.DestinationNodeInfo).Error; err != nil {
						return err
					}
					orderTable.DestinationNodeInfoID = orderTable.DestinationNodeInfo.ID
				}
			}

			if result.OrderTypeID == 0 {
				orderTable.OrderType.OrganizationCountryID = result.OrganizationCountryID
				if err := tx.Create(&orderTable.OrderType).Error; err != nil {
					return err
				}
				orderTable.OrderTypeID = orderTable.OrderType.ID
			}

			orderTable.Commerce = table.Commerce{}
			orderTable.Consumer = table.Consumer{}
			orderTable.OriginContact = table.Contact{}
			orderTable.DestinationContact = table.Contact{}
			orderTable.OriginAddressInfo = table.AddressInfo{}
			orderTable.DestinationAddressInfo = table.AddressInfo{}
			orderTable.OriginNodeInfo = table.NodeInfo{}
			orderTable.DestinationNodeInfo = table.NodeInfo{}
			orderTable.OrderType = table.OrderType{}

			return nil
		})
	}
}
