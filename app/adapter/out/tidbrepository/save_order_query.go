package tidbrepository

import (
	"context"
	views "transport-app/app/adapter/out/tidbrepository/view"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SaveOrderQuery func(ctx context.Context, order domain.Order) (domain.Order, error)

func init() {
	ioc.Registry(
		NewSaveOrderQuery,
		tidb.NewTIDBConnection)
}

func NewSaveOrderQuery(conn tidb.TIDBConnection) SaveOrderQuery {
	return func(ctx context.Context, order domain.Order) (domain.Order, error) {
		var flattenedOrder views.FlattenedOrderView

		err := conn.Raw(`
			SELECT 
				ord.id AS order_id,
				ord.reference_id AS reference_id,
				ord.organization_country_id AS organization_country_id,
				org.country AS organization_country,
				ord.commerce_id AS commerce_id,
				comm.name AS commerce_name,
				ord.consumer_id AS consumer_id,
				cons.name AS consumer_name,
				ord.order_status_id AS order_status_id,
				ost.status AS order_status,
				ord.order_type_id AS order_type_id,
				oty.type AS order_type,
				oty.description AS order_type_description,
				ord.delivery_instructions AS delivery_instructions,
				-- Origen
				ord.origin_contact_id AS origin_contact_id,
				orig_contact.full_name AS origin_contact_name,
				orig_contact.phone AS origin_contact_phone,
				orig_contact.email AS origin_contact_email,
				orig_contact.documents AS origin_contact_documents,
				ord.origin_address_info_id AS origin_address_info_id,
				orig_addr.address_line1 AS origin_address_line1,
				orig_addr.address_line2 AS origin_address_line2,
				orig_addr.address_line3 AS origin_address_line3,
				orig_addr.latitude AS origin_latitude,
				orig_addr.longitude AS origin_longitude,
				orig_addr.state AS origin_state,
				orig_addr.province AS origin_province,
				orig_addr.county AS origin_county,
				orig_addr.district AS origin_district,
				orig_addr.zip_code AS origin_zipcode,
				orig_addr.time_zone AS origin_timezone,
				orig_node.id AS origin_node_info_id,
				orig_node.reference_id AS origin_node_reference_id,
				orig_node.name AS origin_node_name,
				orig_node.type AS origin_node_type,
				orig_node.operator_id AS origin_node_operator_id,
				orig_operator.contact_id AS origin_node_operator_contact_id,

				-- Destino
				ord.destination_contact_id AS destination_contact_id,
				dest_contact.full_name AS destination_contact_name,
				dest_contact.phone AS destination_contact_phone,
				dest_contact.email AS destination_contact_email,
				dest_contact.documents AS destination_contact_documents,
				ord.destination_address_info_id AS destination_address_info_id,
				dest_addr.address_line1 AS destination_address_line1,
				dest_addr.address_line2 AS destination_address_line2,
				dest_addr.address_line3 AS destination_address_line3,
				dest_addr.latitude AS destination_latitude,
				dest_addr.longitude AS destination_longitude,
				dest_addr.state AS destination_state,
				dest_addr.province AS destination_province,
				dest_addr.county AS destination_county,
				dest_addr.district AS destination_district,
				dest_addr.zip_code AS destination_zipcode,
				dest_addr.time_zone AS destination_timezone,
				dest_node.id AS destination_node_info_id,
				dest_node.reference_id AS destination_node_reference_id,
				dest_node.name AS destination_node_name,
				dest_node.type AS destination_node_type,
				dest_node.operator_id AS destination_node_operator_id,
				dest_operator.contact_id AS destination_node_operator_contact_id,

				-- Otros campos
				ord.json_items AS items,
				ord.collect_availability_date AS collect_availability_date,
				ord.collect_availability_time_range_start AS collect_start_time,
				ord.collect_availability_time_range_end AS collect_end_time,
				ord.promised_date_range_start AS promised_start_date,
				ord.promised_date_range_end AS promised_end_date,
				ord.promised_time_range_start AS promised_start_time,
				ord.promised_time_range_end AS promised_end_time,
				ord.transport_requirements AS transport_requirements
			FROM 
				orders ord
				LEFT JOIN organization_countries org ON org.id = ord.organization_country_id
				LEFT JOIN commerces comm ON comm.id = ord.commerce_id AND comm.name = ?
				LEFT JOIN consumers cons ON cons.id = ord.consumer_id AND cons.name = ?
				LEFT JOIN order_statuses ost ON ost.id = ord.order_status_id
				LEFT JOIN order_types oty ON oty.id = ord.order_type_id AND oty.type = ?
				LEFT JOIN contacts orig_contact ON orig_contact.id = ord.origin_contact_id AND orig_contact.full_name = ?
				LEFT JOIN address_infos orig_addr ON orig_addr.id = ord.origin_address_info_id AND orig_addr.raw_address = ?
				LEFT JOIN node_infos orig_node ON orig_node.id = ord.origin_node_info_id AND orig_node.reference_id = ?
				LEFT JOIN operators orig_operator ON orig_operator.id = orig_node.operator_id
				LEFT JOIN contacts dest_contact ON dest_contact.id = ord.destination_contact_id AND dest_contact.full_name = ?
				LEFT JOIN address_infos dest_addr ON dest_addr.id = ord.destination_address_info_id AND dest_addr.raw_address = ?
				LEFT JOIN node_infos dest_node ON dest_node.id = ord.destination_node_info_id AND dest_node.reference_id = ?
				LEFT JOIN operators dest_operator ON dest_operator.id = dest_node.operator_id
			WHERE 
				ord.organization_country_id = ?;
		`,
			order.BusinessIdentifiers.Commerce,
			order.BusinessIdentifiers.Consumer,
			order.OrderType.Type,
			order.Origin.AddressInfo.Contact.FullName,
			order.Origin.AddressInfo.RawAddress(),
			order.Origin.NodeInfo.ReferenceID,
			order.Destination.AddressInfo.Contact.FullName,
			order.Destination.AddressInfo.RawAddress(),
			order.Destination.NodeInfo.ReferenceID,
			order.Organization.OrganizationCountryID,
		).Scan(&flattenedOrder).Error

		var lpns []string
		for _, v := range order.Packages {
			lpns = append(lpns, v.Lpn)
		}

		var flattenedPackages []views.FlattenedPackageView

		err = conn.Raw(`
		SELECT
		    pkg.id AS package_id,
		    pkg.organization_country_id AS organization_country_id,
		    pkg.lpn AS lpn,
		    pkg.json_dimensions->>'$.height' AS height,
		    pkg.json_dimensions->>'$.width' AS width,
		    pkg.json_dimensions->>'$.depth' AS depth,
		    pkg.json_dimensions->>'$.unit' AS unit,
		    pkg.json_weight->>'$.weight_value' AS weight_value,
		    pkg.json_weight->>'$.weight_unit' AS weight_unit,
		    pkg.json_insurance->>'$.unit_value' AS unit_value,
		    pkg.json_insurance->>'$.currency' AS currency
		FROM
		    packages pkg
		WHERE
		    pkg.lpn IN (?) AND
		    pkg.organization_country_id = ?;
			
	`, lpns, order.Organization.OrganizationCountryID).Scan(&flattenedPackages).Error

		return domain.Order{}, err
	}
}
