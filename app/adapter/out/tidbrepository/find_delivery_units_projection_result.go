package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/projection/deliveryunits"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/doug-martin/goqu/v9"
)

type FindDeliveryUnitsProjectionResult func(
	ctx context.Context,
	filters domain.DeliveryUnitsFilter) ([]projectionresult.DeliveryUnitsProjectionResult, error)

func init() {
	ioc.Registry(
		NewFindDeliveryUnitsProjectionResult,
		database.NewConnectionFactory,
		deliveryunits.NewProjection,
	)
}

func NewFindDeliveryUnitsProjectionResult(
	conn database.ConnectionFactory,
	projection deliveryunits.Projection) FindDeliveryUnitsProjectionResult {
	const (
		duh  = "duh"  // delivery_units_histories
		o    = "o"    // orders
		or   = "or"   // order_references
		dadi = "dadi" // destination_address_infos
		oadi = "oadi" // origin_address_infos
		dd   = "dd"   // destination_districts
		dp   = "dp"   // destination_provinces
		dst  = "dst"  // destination_states
		od   = "od"   // origin_districts
		op   = "op"   // origin_provinces
		ost  = "ost"  // origin_states
		du   = "du"   // delivery_units
		dul  = "dul"  // delivery_unit_labels
		oh   = "oh"   // order_headers
		ot   = "ot"   // order_types
		s    = "s"    // status
	)

	return func(ctx context.Context, filters domain.DeliveryUnitsFilter) ([]projectionresult.DeliveryUnitsProjectionResult, error) {
		var results []projectionresult.DeliveryUnitsProjectionResult

		// Dataset base
		ds := goqu.From(goqu.T("delivery_units_histories").As(duh)).
			Select(goqu.I(duh+".id").As("id")).
			InnerJoin(
				goqu.T("orders").As(o),
				goqu.On(goqu.I(o+".document_id").Eq(goqu.I(duh+".order_doc"))),
			).
			Where(goqu.Ex{
				duh + ".tenant_id": sharedcontext.TenantIDFromContext(ctx),
			})

		// Add order references using WITH clause
		if projection.References().Has(filters.RequestedFields) {
			ds = ds.With("order_refs", goqu.From(goqu.T("order_references").As(or)).
				Select(
					goqu.I(or+".order_doc"),
					goqu.L("jsonb_agg(jsonb_build_object('type', type::text, 'value', value::text))").As("references"),
				).
				GroupBy(goqu.I(or+".order_doc")),
			).
				InnerJoin(
					goqu.T("order_refs").As(or),
					goqu.On(goqu.I(or+".order_doc").Eq(goqu.I(o+".document_id"))),
				).
				SelectAppend(goqu.Cast(goqu.I(or+".references"), "jsonb").As("order_references"))
		}

		if projection.DeliveryUnitLabels().Has(filters.RequestedFields) {
			ds = ds.With("delivery_unit_labels", goqu.From(goqu.T("delivery_units_labels").As(dul)).
				Select(
					goqu.I(dul+".delivery_unit_doc"),
					goqu.L("jsonb_agg(jsonb_build_object('type', type::text, 'value', value::text))").As("delivery_unit_labels"),
				).
				GroupBy(goqu.I(dul+".delivery_unit_doc")),
			).
				InnerJoin(
					goqu.T("delivery_unit_labels").As(dul),
					goqu.On(goqu.I(dul+".delivery_unit_doc").Eq(goqu.I(duh+".delivery_unit_doc"))),
				).
				SelectAppend(goqu.Cast(goqu.I(dul+".delivery_unit_labels"), "jsonb").As("delivery_unit_labels"))
		}

		// Campos de delivery_units_histories
		if projection.Channel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".channel").As("channel"))
		}

		if projection.Status().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("statuses").As(s),
				goqu.On(goqu.I(s+".document_id").Eq(goqu.I(duh+".delivery_unit_status_doc"))),
			).
				SelectAppend(goqu.I(s + ".status").As("status"))
		}

		// LPN and Package Information
		if projection.DeliveryUnit().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("delivery_units").As(du),
				goqu.On(goqu.I(du+".document_id").Eq(goqu.I(duh+".delivery_unit_doc"))),
			)
		}

		if projection.DeliveryUnitLPN().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".lpn").As("lpn"))
		}

		if projection.DeliveryUnitDimensions().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_dimensions").As("json_dimensions"))
		}
		if projection.DeliveryUnitWeight().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_weight").As("json_weight"))
		}
		if projection.DeliveryUnitInsurance().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_insurance").As("json_insurance"))
		}
		if projection.DeliveryUnitItems().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_items").As("json_items"))
		}

		// Campos de orders
		if projection.ReferenceID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".reference_id").As("order_reference_id"))
		}

		if projection.GroupByType().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".group_by_type").As("order_group_by_type"))
		}

		if projection.GroupByValue().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".group_by_value").As("order_group_by_value"))
		}

		if projection.CollectAvailabilityDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_date").As("order_collect_availability_date"))
		}

		if projection.DeliveryInstructions().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".delivery_instructions").As("order_delivery_instructions"))
		}

		if projection.CollectAvailabilityDateStartTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_time_range_start").As("order_collect_availability_date_start_time"))
		}
		if projection.CollectAvailabilityDateEndTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_time_range_end").As("order_collect_availability_date_end_time"))
		}

		if projection.PromisedDateDateRangeStartDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_date_range_start").As("order_promised_date_start_date"))
		}
		if projection.PromisedDateDateRangeEndDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_date_range_end").As("order_promised_date_end_date"))
		}
		if projection.PromisedDateTimeRangeStartTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_time_range_start").As("order_promised_date_start_time"))
		}
		if projection.PromisedDateTimeRangeEndTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_time_range_end").As("order_promised_date_end_time"))
		}
		if projection.PromisedDateServiceCategory().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".service_category").As("order_promised_date_service_category"))
		}

		// Campos de orderType
		if projection.OrderType().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("order_types").As(ot),
				goqu.On(goqu.I(ot+".document_id").Eq(goqu.I(o+".order_type_doc"))),
			)
		}

		if projection.OrderTypeType().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ot + ".type").As("order_type"))
		}

		if projection.OrderTypeDescription().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ot + ".description").As("order_type_description"))
		}

		// Campos de address_infos
		if projection.DestinationAddressLine2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".address_line2").As("destination_address_line2"))
		}

		if projection.Commerce().Has(filters.RequestedFields) || projection.Consumer().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("order_headers").As(oh),
				goqu.On(goqu.I(oh+".document_id").Eq(goqu.I(o+".order_headers_doc"))),
			)
		}

		if projection.Commerce().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oh + ".commerce").As("commerce"))
		}

		if projection.Consumer().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oh + ".consumer").As("consumer"))
		}

		// Join address_infos si se requiere algún campo de addressInfo
		if projection.DestinationAddressInfo().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("address_infos").As(dadi),
				goqu.On(goqu.I(dadi+".document_id").Eq(goqu.I(o+".destination_address_info_doc"))),
			)
		}

		// Campos de address_infos
		if projection.DestinationAddressLine1().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".address_line1").As("destination_address_line1"))
		}

		// Join con districts
		if projection.DestinationDistrict().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("districts").As(dd),
				goqu.On(goqu.I(dd+".document_id").Eq(goqu.I(dadi+".district_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dd + ".name").As("destination_district"))
		}

		// Join con provinces
		if projection.DestinationProvince().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("provinces").As(dp),
				goqu.On(goqu.I(dp+".document_id").Eq(goqu.I(dadi+".province_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dp + ".name").As("destination_province"))
		}

		// Join con states
		if projection.DestinationState().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("states").As(dst),
				goqu.On(goqu.I(dst+".document_id").Eq(goqu.I(dadi+".state_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dst + ".name").As("destination_state"))
		}

		if projection.DestinationZipCode().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".zip_code").As("destination_zip_code"))
		}

		// Campos de coordenadas
		if projection.DestinationCoordinatesLatitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".latitude").As("destination_coordinates_latitude"))
		}

		if projection.DestinationCoordinatesLongitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".longitude").As("destination_coordinates_longitude"))
		}

		if projection.DestinationCoordinatesSource().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_source").As("destination_coordinates_source"))
		}

		if projection.DestinationCoordinatesConfidenceLevel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_confidence").As("destination_coordinates_confidence_level"))
		}

		if projection.DestinationCoordinatesConfidenceMessage().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_message").As("destination_coordinates_confidence_message"))
		}

		if projection.DestinationCoordinatesConfidenceReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_reason").As("destination_coordinates_confidence_reason"))
		}

		if projection.DestinationTimeZone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".time_zone").As("destination_time_zone"))
		}

		// Campos de contacto del destino
		if projection.DestinationContact().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("contacts").As("dc"),
				goqu.On(goqu.I("dc.document_id").Eq(goqu.I(o+".destination_contact_doc"))),
			)
		}

		if projection.DestinationContactEmail().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.email").As("destination_contact_email"))
		}

		if projection.DestinationContactFullName().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.full_name").As("destination_contact_full_name"))
		}

		if projection.DestinationContactNationalID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.national_id").As("destination_contact_national_id"))
		}

		if projection.DestinationContactPhone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.phone").As("destination_contact_phone"))
		}

		if projection.DestinationAdditionalContactMethods().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.additional_contact_methods").As("destination_additional_contact_methods"))
		}

		if projection.DestinationContactDocuments().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.documents").As("destination_contact_documents"))
		}

		// Join address_infos para origen si se requiere algún campo de addressInfo
		if projection.OriginAddressInfo().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("address_infos").As(oadi),
				goqu.On(goqu.I(oadi+".document_id").Eq(goqu.I(o+".origin_address_info_doc"))),
			)
		}

		// Campos de address_infos para origen
		if projection.OriginAddressLine1().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".address_line1").As("origin_address_line1"))
		}

		if projection.OriginAddressLine2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".address_line2").As("origin_address_line2"))
		}

		// Join con districts para origen
		if projection.OriginDistrict().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("districts").As(od),
				goqu.On(goqu.I(od+".document_id").Eq(goqu.I(oadi+".district_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(od + ".name").As("origin_district"))
		}

		// Join con provinces para origen
		if projection.OriginProvince().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("provinces").As(op),
				goqu.On(goqu.I(op+".document_id").Eq(goqu.I(oadi+".province_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(op + ".name").As("origin_province"))
		}

		// Join con states para origen
		if projection.OriginState().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("states").As(ost),
				goqu.On(goqu.I(ost+".document_id").Eq(goqu.I(oadi+".state_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(ost + ".name").As("origin_state"))
		}

		if projection.OriginZipCode().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".zip_code").As("origin_zip_code"))
		}

		// Campos de coordenadas para origen
		if projection.OriginCoordinatesLatitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".latitude").As("origin_coordinates_latitude"))
		}

		if projection.OriginCoordinatesLongitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".longitude").As("origin_coordinates_longitude"))
		}

		if projection.OriginCoordinatesSource().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_source").As("origin_coordinates_source"))
		}

		if projection.OriginCoordinatesConfidenceLevel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_confidence").As("origin_coordinates_confidence_level"))
		}

		if projection.OriginCoordinatesConfidenceMessage().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_message").As("origin_coordinates_confidence_message"))
		}

		if projection.OriginCoordinatesConfidenceReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_reason").As("origin_coordinates_confidence_reason"))
		}

		if projection.OriginTimeZone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".time_zone").As("origin_time_zone"))
		}

		// Campos de contacto del origen
		if projection.OriginContact().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("contacts").As("oc"),
				goqu.On(goqu.I("oc.document_id").Eq(goqu.I(o+".origin_contact_doc"))),
			)
		}

		if projection.OriginContactEmail().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.email").As("origin_contact_email"))
		}

		if projection.OriginContactFullName().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.full_name").As("origin_contact_full_name"))
		}

		if projection.OriginContactNationalID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.national_id").As("origin_contact_national_id"))
		}

		if projection.OriginContactPhone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.phone").As("origin_contact_phone"))
		}

		if projection.OriginContactMethods().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.additional_contact_methods").As("origin_additional_contact_methods"))
		}

		if projection.OriginDocuments().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.documents").As("origin_contact_documents"))
		}

		if projection.ExtraFields().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".extra_fields").As("extra_fields"))
		}

		sql, args, err := ds.Prepared(true).ToSQL()
		if err != nil {
			return nil, err
		}

		err = conn.WithContext(ctx).Raw(sql, args...).Scan(&results).Error
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}
