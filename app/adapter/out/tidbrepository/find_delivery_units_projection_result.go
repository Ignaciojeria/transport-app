package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
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
		database.NewConnectionFactory)
}

func NewFindDeliveryUnitsProjectionResult(
	conn database.ConnectionFactory) FindDeliveryUnitsProjectionResult {
	const (
		duh = "duh"
		o   = "o"
	)
	return func(ctx context.Context, filters domain.DeliveryUnitsFilter) ([]projectionresult.DeliveryUnitsProjectionResult, error) {

		var results []projectionresult.DeliveryUnitsProjectionResult

		deliveryUnitHistoryIndex := []interface{}{
			goqu.I(duh + ".id").As("id"),
		}

		ordersSelectProjection := []interface{}{
			goqu.I(o + ".reference_id").As("order_reference_id"),
		}

		ds := goqu.From(goqu.T("delivery_units_histories").As(duh)).
			Select(deliveryUnitHistoryIndex...).
			SelectAppend(ordersSelectProjection...).
			InnerJoin(
				goqu.T("orders").As(o),
				goqu.On(goqu.I(o+".document_id").Eq(goqu.I(duh+".order_doc"))),
			).
			Where(goqu.Ex{
				duh + ".tenant_id": sharedcontext.TenantIDFromContext(ctx),
			})

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
