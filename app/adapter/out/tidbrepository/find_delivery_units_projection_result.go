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
		dadi = "dadi" // destination_address_infos
		dd   = "dd"   // destination_districts
		dp   = "dp"   // destination_provinces
		dst  = "dst"  // destination_states
		du   = "du"   // delivery_units
		oh   = "oh"   // order_headers
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

		// Campos de delivery_units_histories
		if projection.Channel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".channel").As("channel"))
		}

		// LPN and Package Information
		if projection.DeliveryUnit().Has(filters.RequestedFields) {
			ds = ds.LeftJoin(
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

		if projection.CollectAvailabilityDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_date").As("order_collect_availability_date"))
		}
		/*
			if projection.DestinationDeliveryInstructions().Has(filters.RequestedFields) {
				ds = ds.SelectAppend(goqu.I(o + ".delivery_instructions").As("order_delivery_instructions"))
			}
		*/
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

		// Campos de address_infos
		if projection.DestinationAddressLine2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".address_line2").As("destination_address_line2"))
		}

		if projection.Commerce().Has(filters.RequestedFields) || projection.Consumer().Has(filters.RequestedFields) {
			ds = ds.LeftJoin(
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
			ds = ds.LeftJoin(
				goqu.T("districts").As(dd),
				goqu.On(goqu.I(dd+".document_id").Eq(goqu.I(dadi+".district_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dd + ".name").As("destination_district"))
		}

		// Join con provinces
		if projection.DestinationProvince().Has(filters.RequestedFields) {
			ds = ds.LeftJoin(
				goqu.T("provinces").As(dp),
				goqu.On(goqu.I(dp+".document_id").Eq(goqu.I(dadi+".province_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dp + ".name").As("destination_province"))
		}

		// Join con states
		if projection.DestinationState().Has(filters.RequestedFields) {
			ds = ds.LeftJoin(
				goqu.T("states").As(dst),
				goqu.On(goqu.I(dst+".document_id").Eq(goqu.I(dadi+".state_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dst + ".name").As("destination_state"))
		}

		if projection.DestinationLatitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".latitude").As("destination_latitude"))
		}
		if projection.DestinationLongitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".longitude").As("destination_longitude"))
		}
		if projection.DestinationTimeZone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".time_zone").As("destination_time_zone"))
		}

		if projection.DestinationZipCode().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".zip_code").As("destination_zip_code"))
		}

		if projection.DestinationRequiresManualReview().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".requires_manual_review").As("destination_requires_manual_review"))
		}

		// Campos de contacto del destino
		if projection.DestinationContact().Has(filters.RequestedFields) {
			ds = ds.LeftJoin(
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
