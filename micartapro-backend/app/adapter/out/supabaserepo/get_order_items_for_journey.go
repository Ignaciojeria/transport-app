package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"

	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

// OrderItemRow es una fila de order_items_projection para el reporte.
type OrderItemRow struct {
	AggregateID   int64   `json:"aggregate_id"`
	OrderNumber   *int    `json:"order_number"`
	ItemName      string  `json:"item_name"`
	Quantity      int     `json:"quantity"`
	Unit          string  `json:"unit"`
	Station       *string `json:"station"`
	Fulfillment   string  `json:"fulfillment"`
	Status        string  `json:"status"`
	RequestedTime *string `json:"requested_time"`
	CreatedAt     string  `json:"created_at"`
}

// GetOrderItemsForJourney obtiene todos los ítems de order_items_projection para una jornada.
type GetOrderItemsForJourney func(ctx context.Context, journeyID string) ([]OrderItemRow, error)

func init() {
	ioc.Registry(NewGetOrderItemsForJourney, supabasecli.NewSupabaseClient)
}

func NewGetOrderItemsForJourney(supabase *supabase.Client) GetOrderItemsForJourney {
	return func(ctx context.Context, journeyID string) ([]OrderItemRow, error) {
		data, _, err := supabase.From("order_items_projection").
			Select("aggregate_id,order_number,item_name,quantity,unit,station,fulfillment,status,requested_time,created_at", "", false).
			Eq("journey_id", journeyID).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("querying order_items_projection: %w", err)
		}
		var rows []OrderItemRow
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("unmarshaling order_items_projection: %w", err)
		}
		return rows, nil
	}
}
