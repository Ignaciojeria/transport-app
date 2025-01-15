package tidbrepository

import (
	"context"
	"fmt"
	views "transport-app/app/adapter/out/tidbrepository/view"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewFindOrdersByReferenceAndFilters,
		tidb.NewTIDBConnection)
}

type FindOrdersByReferenceAndFilters func(
	context.Context,
	domain.OrderSearchFilters) ([]domain.Order, error)

func NewFindOrdersByReferenceAndFilters(conn tidb.TIDBConnection) FindOrdersByReferenceAndFilters {
	return func(ctx context.Context, osf domain.OrderSearchFilters) ([]domain.Order, error) {
		var orders []views.FlattenedOrderView

		// Query principal para obtener las órdenes
		query := `
SELECT 
    o.id as order_id,
    o.reference_id,
    org_country.country as organization_country,
    com.name as commerce_name,
    con.name as consumer_name,
    os.status as order_status,
    ot.type as order_type,
    ot.description as order_type_description,
    o.delivery_instructions,
    -- Datos del contacto de origen
    oc.full_name as origin_contact_name,
    oc.phone as origin_contact_phone,
    oc.email as origin_contact_email,
    oc.national_id as origin_contact_national_id, -- Suponiendo que este campo existe
    oc.documents as origin_contact_documents,
    -- Datos del contacto de destino
    dc.full_name as destination_contact_name,
    dc.phone as destination_contact_phone,
    dc.email as destination_contact_email,
    dc.national_id as destination_contact_national_id, -- Suponiendo que este campo existe
    dc.documents as destination_contact_documents,
    -- Datos de dirección de origen
    oa.address_line1 as origin_address_line1,
    oa.address_line2 as origin_address_line2,
    oa.address_line3 as origin_address_line3,
    oa.state as origin_state,
    oa.province as origin_province,
    oa.county as origin_county,
    oa.district as origin_district,
    oa.zip_code as origin_zipcode,
    oa.latitude as origin_latitude,
    oa.longitude as origin_longitude,
    oa.time_zone as origin_timezone,
    -- Datos del nodo de origen
    on_info.reference_id as origin_node_reference_id,
    on_info.name as origin_node_name,
    on_info.type as origin_node_type,
    on_op.full_name as origin_node_operator_name,
    -- Datos de dirección de destino
    da.address_line1 as destination_address_line1,
    da.address_line2 as destination_address_line2,
    da.address_line3 as destination_address_line3,
    da.state as destination_state,
    da.province as destination_province,
    da.county as destination_county,
    da.district as destination_district,
    da.zip_code as destination_zipcode,
    da.latitude as destination_latitude,
    da.longitude as destination_longitude,
    da.time_zone as destination_timezone,
    -- Datos del nodo de destino
    dn_info.reference_id as destination_node_reference_id,
    dn_info.name as destination_node_name,
    dn_info.type as destination_node_type,
    dn_op.full_name as destination_node_operator_name,
    -- Otros datos
    o.json_items as items,
    o.collect_availability_date,
    o.collect_availability_time_range_start as collect_start_time,
    o.collect_availability_time_range_end as collect_end_time,
    o.promised_date_range_start as promised_start_date,
    o.promised_date_range_end as promised_end_date,
    o.promised_time_range_start as promised_start_time,
    o.promised_time_range_end as promised_end_time,
    o.transport_requirements
FROM orders o
LEFT JOIN commerces com ON o.commerce_id = com.id
LEFT JOIN consumers con ON o.consumer_id = con.id
LEFT JOIN order_statuses os ON o.order_status_id = os.id
LEFT JOIN order_types ot ON o.order_type_id = ot.id
LEFT JOIN organization_countries org_country ON o.organization_country_id = org_country.id
LEFT JOIN contacts oc ON o.origin_contact_id = oc.id
LEFT JOIN address_infos oa ON o.origin_address_info_id = oa.id
LEFT JOIN node_infos on_info ON o.origin_node_info_id = on_info.id
LEFT JOIN contacts on_op ON on_info.operator_id = on_op.id
LEFT JOIN contacts dc ON o.destination_contact_id = dc.id
LEFT JOIN address_infos da ON o.destination_address_info_id = da.id
LEFT JOIN node_infos dn_info ON o.destination_node_info_id = dn_info.id
LEFT JOIN contacts dn_op ON dn_info.operator_id = dn_op.id
WHERE 
    o.reference_id IN (?) 
    AND org_country.country = ? 
    AND EXISTS (
        SELECT 1 
        FROM api_keys ak
        JOIN organizations org ON ak.organization_id = org.id
        JOIN organization_countries org_country_filter ON org.id = org_country_filter.organization_id
        WHERE ak.key = ? 
        AND org_country_filter.country = ?
    )
    `

		params := []interface{}{
			osf.ReferenceIDs,
			osf.Organization.Country.Alpha2(),
			osf.Organization.Key,
			osf.Organization.Country.Alpha2(),
		}

		if err := conn.Raw(query, params...).Scan(&orders).Error; err != nil {
			return nil, fmt.Errorf("error scanning orders: %w", err)
		}

		if len(orders) == 0 {
			return []domain.Order{}, nil
		}

		// Obtener IDs de órdenes
		orderIDs := make([]int64, len(orders))
		for i, order := range orders {
			orderIDs[i] = order.OrderID
		}

		// Obtener packages
		var packages []views.FlattenedPackageView
		packagesQuery := `
        SELECT 
            op.order_id,
            p.lpn,
            p.json_dimensions->>'$.Height' as height,
            p.json_dimensions->>'$.Width' as width,
            p.json_dimensions->>'$.Depth' as depth,
            p.json_dimensions->>'$.Unit' as unit,
            p.json_weight->>'$.WeightValue' as weight_value,
            p.json_weight->>'$.WeightUnit' as weight_unit,
            p.json_items->>'$.Description' as description,
            p.json_items->>'$.QuantityNumber' as quantity,
            p.json_insurance->>'$.UnitValue' as unit_value,
            p.json_insurance->>'$.Currency' as currency,
            'default' as package_type
        FROM packages p
        JOIN order_packages op ON p.id = op.package_id
        WHERE op.order_id IN (?)`

		if err := conn.Raw(packagesQuery, orderIDs).Scan(&packages).Error; err != nil {
			return nil, fmt.Errorf("error scanning packages: %w", err)
		}

		// Obtener referencias
		var references []views.FlattenedOrderReferenceView
		referencesQuery := `
        SELECT 
            order_id,
            type,
            value
        FROM order_references
        WHERE order_id IN (?)`

		if err := conn.Raw(referencesQuery, orderIDs).Scan(&references).Error; err != nil {
			return nil, fmt.Errorf("error scanning references: %w", err)
		}

		// Obtener visits
		var visits []views.FlattenedVisitView
		visitsQuery := `
        SELECT 
            order_id,
            date,
            time_range_start,
            time_range_end
        FROM visits
        WHERE order_id IN (?)
        ORDER BY date`

		if err := conn.Raw(visitsQuery, orderIDs).Scan(&visits).Error; err != nil {
			return nil, fmt.Errorf("error scanning visits: %w", err)
		}

		// Agrupar por OrderID
		packagesByOrder := make(map[int64][]views.FlattenedPackageView)
		for _, p := range packages {
			packagesByOrder[p.OrderID] = append(packagesByOrder[p.OrderID], p)
		}

		referencesByOrder := make(map[int64][]views.FlattenedOrderReferenceView)
		for _, r := range references {
			referencesByOrder[r.OrderID] = append(referencesByOrder[r.OrderID], r)
		}

		visitsByOrder := make(map[int64][]views.FlattenedVisitView)
		for _, v := range visits {
			visitsByOrder[v.OrderID] = append(visitsByOrder[v.OrderID], v)
		}

		// Mapear a domain.Order
		domainOrders := make([]domain.Order, len(orders))
		for i, order := range orders {
			domainOrders[i] = order.ToOrder(
				packagesByOrder[order.OrderID],
				referencesByOrder[order.OrderID],
				visitsByOrder[order.OrderID],
			)
		}

		return domainOrders, nil
	}
}
