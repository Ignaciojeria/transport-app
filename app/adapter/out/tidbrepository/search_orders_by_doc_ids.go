package tidbrepository

//TODO IMPLEMENT THIS QUERY :
/*
SELECT
  o.id,
  o.reference_id as order_reference_id,
  ot.type as order_type,
  ot.description as order_type_description,
  o.delivery_instructions as order_delivery_instructions,
  o.collect_availability_date as order_collect_availability_date,
  o.collect_availability_time_range_start as order_collect_availability_time_range_start,
  o.collect_availability_time_range_end as order_collect_availability_time_range_end,
  o.service_category as order_service_category,
  o.promised_date_range_start as order_promised_date_range_start,
  o.promised_date_range_end as order_promised_date_range_end,
  o.promised_time_range_start as order_promised_time_range_start,
  o.promised_time_range_end as order_promised_time_range_end,
  o.transport_requirements as order_transport_requirements,
  o.address_line2 as order_address_line2,
  oh.commerce as order_headers_commerce,
  oh.consumer as order_headers_consumer,
  oc.documents as origin_contact_documents,
  oc.email as origin_contact_email,
  oc.full_name as origin_contact_fullname,
  oc.phone as origin_contact_phone,
  oc.national_id as origin_contact_national_id,
  dc.email as destination_contact_email,
  dc.full_name as destination_contact_name,
  dc.documents as destination_contact_documents,
  dc.national_id as destination_contact_national_id,
  dc.phone as destination_contact_phone,
  pk.lpn as package_lpn,
  pk.json_dimensions as package_json_dimensions,
  pk.json_weight as package_json_weight,
  pk.json_insurance as package_json_insurance,
  pk.json_items as package_json_items,
  oni.name as origin_node_name,
  oni.reference_id as origin_node_id,
  dni.reference_id as destination_node_id,
  dni.name as node_name,
  oadi.state as origin_address_state,
  oadi.province as origin_address_province,
  oadi.district as origin_address_district,
  oadi.address_line1 as origin_address_line1,
  oadi.latitude as origin_address_latitude,
  oadi.longitude as origin_address_longitude,
  oadi.time_zone as origin_time_zone,
  oadi.zip_code as origin_zip_code,
  dadi.state as destination_address_state,
  dadi.province as destination_address_province,
  dadi.district as destination_address_district,
  dadi.address_line1 as destination_address_line1,
  dadi.latitude as destination_address_latitude,
  dadi.longitude as destination_address_longitude,
  dadi.time_zone as destination_time_zone,
  dadi.zip_code as destination_zip_code,
  (
    SELECT JSON_ARRAYAGG(JSON_OBJECT('type', type, 'value', value))
    FROM order_references
    WHERE order_doc = o.document_id
  ) AS order_references
FROM orders o
INNER JOIN order_headers oh ON o.order_headers_doc = oh.document_id
INNER JOIN contacts oc ON oc.document_id = o.origin_contact_doc
INNER JOIN contacts dc ON dc.document_id = o.destination_contact_doc
INNER JOIN order_packages op on op.order_doc = o.document_id
INNER JOIN packages pk on pk.document_id = op.package_doc
INNER JOIN node_infos oni on o.origin_node_info_doc = oni.document_id
INNER JOIN node_infos dni on dni.document_id = o.destination_node_info_doc
INNER JOIN address_infos oadi on oadi.document_id = o.origin_address_info_doc
INNER JOIN address_infos dadi on dadi.document_id = o.destination_address_info_doc
INNER JOIN order_types ot on ot.document_id = o.order_type_doc
WHERE o.document_id IN ('a7c1bd4415413d12d2ffc2abd2d14f7e')
AND o.organization_id = 1;
*/
