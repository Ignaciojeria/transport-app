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
		// Campos de orders
		if projection.ReferenceID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".reference_id").As("order_reference_id"))
		}

		if projection.CollectAvailabilityDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_date").As("order_collect_availability_date"))
		}
		if projection.CollectAvailabilityDateStartTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_time_range_start").As("order_collect_availability_date_start_time"))
		}
		if projection.CollectAvailabilityDateEndTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_time_range_end").As("order_collect_availability_date_end_time"))
		}

		// Campos de PromisedDate
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

		// Join address_infos si se requiere algún campo de addressInfo
		if projection.DestinationAddressInfo().HasAnyPrefix(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("address_infos").As(dadi),
				goqu.On(goqu.I(dadi+".document_id").Eq(goqu.I(o+".destination_address_info_doc"))),
			)
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
