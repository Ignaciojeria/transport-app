package tidbrepository

import (
	"context"
	"fmt"
	views "transport-app/app/adapter/out/tidbrepository/views"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewFindOrdersByFilters,
		database.NewConnectionFactory)
}

type FindOrdersByFilters func(
	context.Context,
	domain.OrderSearchFilters) ([]domain.Order, error)

func NewFindOrdersByFilters(conn database.ConnectionFactory) FindOrdersByFilters {
	return func(ctx context.Context, osf domain.OrderSearchFilters) ([]domain.Order, error) {
		var orders []views.FlattenedOrderView

		// Query principal para obtener las órdenes
		query := `
        SELECT DISTINCT
            o.id as order_id,
            o.reference_id,
            o.sequence_number as order_sequence_number,
            headers.commerce as commerce_name,
            headers.consumer as consumer_name,
            os.status as order_status,
            ot.type as order_type,
            ot.description as order_type_description,
            o.delivery_instructions,
            -- Datos del contacto de origen
            oc.full_name as origin_contact_name,
            oc.phone as origin_contact_phone,
            oc.email as origin_contact_email,
            oc.national_id as origin_contact_national_id,
            oc.documents as origin_contact_documents,
            -- Datos del contacto de destino
            dc.full_name as destination_contact_name,
            dc.phone as destination_contact_phone,
            dc.email as destination_contact_email,
            dc.national_id as destination_contact_national_id,
            dc.documents as destination_contact_documents,
            -- Datos de dirección de origen
            oa.address_line1 as origin_address_line1,
            oa.address_line2 as origin_address_line2,
            oa.address_line3 as origin_address_line3,
            oa.state as origin_state,
            oa.province as origin_province,
            oa.locality as origin_locality,
            oa.district as origin_district,
            oa.zip_code as origin_zipcode,
            oa.latitude as origin_latitude,
            oa.longitude as origin_longitude,
            oa.time_zone as origin_timezone,
            -- Datos del nodo de origen
            on_info.reference_id as origin_node_reference_id,
            on_info.name as origin_node_name,
            -- Datos de dirección de destino
            da.address_line1 as destination_address_line1,
            da.address_line2 as destination_address_line2,
            da.address_line3 as destination_address_line3,
            da.state as destination_state,
            da.province as destination_province,
            da.locality as destination_locality,
            da.district as destination_district,
            da.zip_code as destination_zipcode,
            da.latitude as destination_latitude,
            da.longitude as destination_longitude,
            da.time_zone as destination_timezone,
            -- Datos del nodo de destino
            dn_info.reference_id as destination_node_reference_id,
            dn_info.name as destination_node_name,
            -- Otros datos
            o.json_items as items,
            o.collect_availability_date,
            o.collect_availability_time_range_start as collect_start_time,
            o.collect_availability_time_range_end as collect_end_time,
            o.promised_date_range_start as promised_start_date,
            o.promised_date_range_end as promised_end_date,
            o.promised_time_range_start as promised_start_time,
            o.promised_time_range_end as promised_end_time,
            o.transport_requirements,
            -- Campos nuevos relacionados con el Plan
            o.plan_id,
            pln.planned_date,
            pln.reference_id AS plan_reference_id,
            pln.json_start_location as plan_start_location,
            -- Campos nuevos relacionados con la Ruta
            r.id AS route_id,
            r.reference_id AS route_reference_id,
            r.json_end_location AS route_end_location,
            r.end_node_reference_id AS route_end_node_reference_id,
            r.account_id AS route_account_id,
            a.reference_id AS route_account_reference_id
        FROM orders o
        LEFT JOIN order_headers headers ON o.order_headers_id = headers.id
        LEFT JOIN order_statuses os ON o.order_status_id = os.id
        LEFT JOIN order_types ot ON o.order_type_id = ot.id
        LEFT JOIN organizations org ON o.organization_id = org.id  -- Se mantiene el filtro por Organization ID
        LEFT JOIN contacts oc ON o.origin_contact_id = oc.id
        LEFT JOIN address_infos oa ON o.origin_address_info_id = oa.id
        LEFT JOIN node_infos on_info ON o.origin_node_info_id = on_info.id
        LEFT JOIN contacts dc ON o.destination_contact_id = dc.id
        LEFT JOIN address_infos da ON o.destination_address_info_id = da.id
        LEFT JOIN node_infos dn_info ON o.destination_node_info_id = dn_info.id
        LEFT JOIN plans pln ON o.plan_id = pln.id
        LEFT JOIN routes r ON o.route_id = r.id
        LEFT JOIN accounts a ON r.account_id = a.id
        LEFT JOIN order_packages op ON o.id = op.order_id
        LEFT JOIN packages p ON op.package_id = p.id
        WHERE 
            o.organization_id = ?`

		params := []interface{}{
			sharedcontext.TenantIDFromContext(ctx),
		}

		// Condición dinámica para LPNs
		if len(osf.Lpns) > 0 {
			query += " AND p.lpn IN (?)"
			params = append(params, osf.Lpns)
		}

		// Condición dinámica para ReferenceIDs
		if len(osf.ReferenceIDs) > 0 {
			query += " AND o.reference_id IN (?)"
			params = append(params, osf.ReferenceIDs)
		}

		// Condición dinámica para OperatorReferenceID

		if osf.PlanReferenceID != "" {
			query += " AND pln.reference_id = ?"
			params = append(params, osf.PlanReferenceID)
		}

		// Condición adicional para comercios
		if len(osf.Commerces) > 0 {
			query += " AND headers.commerce IN (?)"
			params = append(params, osf.Commerces)
		}

		// Ejecutar la consulta
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
                p.json_dimensions->>'$.height' AS height,
                p.json_dimensions->>'$.width' AS width,
                p.json_dimensions->>'$.Length' AS Length,
                p.json_dimensions->>'$.unit' AS unit,
                p.json_weight->>'$.weight_value' AS weight_value,
                p.json_weight->>'$.weight_unit' AS weight_unit,
                p.json_items_references as items_references,
                p.json_insurance->>'$.unit_value' AS unit_value,
                p.json_insurance->>'$.currency' AS currency,
                'default' AS package_type
            FROM packages p
            JOIN order_packages op ON p.id = op.package_id
            WHERE op.order_id IN (?);
            `

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

		// Agrupar por OrderID
		packagesByOrder := make(map[int64][]views.FlattenedPackageView)
		for _, p := range packages {
			packagesByOrder[p.OrderID] = append(packagesByOrder[p.OrderID], p)
		}

		referencesByOrder := make(map[int64][]views.FlattenedOrderReferenceView)
		for _, r := range references {
			referencesByOrder[r.OrderID] = append(referencesByOrder[r.OrderID], r)
		}

		// Mapear a domain.Order
		domainOrders := make([]domain.Order, len(orders))
		for i, order := range orders {
			domainOrders[i] = order.ToOrder(
				packagesByOrder[order.OrderID],
				referencesByOrder[order.OrderID],
			)
		}

		return domainOrders, nil
	}
}
