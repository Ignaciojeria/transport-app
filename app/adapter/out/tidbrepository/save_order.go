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

		var QueryResult struct {
			OrganizationCountryID int64
			CommerceID            int64
			ConsumerID            int64
			OrderTypeID           int64
			ContactID             int64
			AddressID             int64
			OriginNodeInfoID      int64
			DestinationNodeInfoID int64
		}

		err := conn.Raw(`
			WITH 
			  org_country AS (
				SELECT 
				  oc.id AS organization_country_id
				FROM 
				  organization_countries oc
				INNER JOIN organizations o ON o.id = oc.organization_id
				WHERE 
				  oc.country = ? AND o.key = ?
			  ),
			  commerce AS (
				SELECT 
				  c.id AS commerce_id
				FROM 
				  commerces c
				WHERE 
				  c.organization_country_id = (SELECT organization_country_id FROM org_country)
				  AND c.name = ?
			  ),
			  consumer AS (
				SELECT 
				  con.id AS consumer_id
				FROM 
				  consumers con
				WHERE 
				  con.organization_country_id = (SELECT organization_country_id FROM org_country)
				  AND con.name = ?
			  ),
			  order_type AS (
				SELECT 
				  ot.id AS order_type_id
				FROM 
				  order_types ot
				WHERE 
				  ot.organization_country_id = (SELECT organization_country_id FROM org_country)
				  AND ot.type = ?
			  ),
			  contact AS (
				SELECT 
				  ct.id AS contact_id
				FROM 
				  contacts ct
				WHERE 
				  ct.organization_country_id = (SELECT organization_country_id FROM org_country)
				  AND ct.full_name = ?
				  AND ct.email = ?
				  AND ct.phone = ?
			  ),
			  address AS (
				SELECT 
				  ai.id AS address_id
				FROM 
				  address_infos ai
				WHERE 
				  ai.organization_country_id = (SELECT organization_country_id FROM org_country)
				  AND ai.raw_address = ?
			  ),
			  origin_node AS (
				SELECT 
				  ni.id AS origin_node_info_id
				FROM 
				  node_infos ni
				WHERE 
				  ni.organization_country_id = (SELECT organization_country_id FROM org_country)
				  AND ni.reference_id = ?
			  ),
			  destination_node AS (
				SELECT 
				  ni.id AS destination_node_info_id
				FROM 
				  node_infos ni
				WHERE 
				  ni.organization_country_id = (SELECT organization_country_id FROM org_country)
				  AND ni.reference_id = ?
			  )
			SELECT 
			  (SELECT organization_country_id FROM org_country) AS organization_country_id,
			  (SELECT commerce_id FROM commerce) AS commerce_id,
			  (SELECT consumer_id FROM consumer) AS consumer_id,
			  (SELECT order_type_id FROM order_type) AS order_type_id,
			  (SELECT contact_id FROM contact) AS contact_id,
			  (SELECT address_id FROM address) AS address_id,
			  (SELECT origin_node_info_id FROM origin_node) AS origin_node_info_id,
			  (SELECT destination_node_info_id FROM destination_node) AS destination_node_info_id;
		`,
			to.Organization.Country,
			to.Organization.Key,
			to.BusinessIdentifiers.Commerce,
			to.BusinessIdentifiers.Consumer,
			to.OrderType.Type,
			to.Origin.AddressInfo.Contact.FullName,
			to.Origin.AddressInfo.Contact.Email,
			to.Origin.AddressInfo.Contact.Phone,
			to.Origin.AddressInfo.RawAddress,
			to.Origin.NodeInfo.ReferenceID,
			to.Destination.NodeInfo.ReferenceID,
		).Scan(&QueryResult).Error

		if err != nil {
			return domain.Order{}, err
		}

		return domain.Order{}, err
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
