package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// OrderItemByAggregateRow es una fila de order_items_projection para una orden por aggregate_id.
type OrderItemByAggregateRow struct {
	ItemName   string  `json:"item_name"`
	Quantity   int     `json:"quantity"`
	Unit       string  `json:"unit"`
	TotalPrice float64 `json:"total_price"`
	Station    *string `json:"station"`
}

// GetOrderItemsByAggregateID obtiene los ítems de order_items_projection para un aggregate_id y menu_id.
type GetOrderItemsByAggregateID func(ctx context.Context, menuID string, aggregateID int64) ([]OrderItemByAggregateRow, error)

func init() {
	ioc.Register(NewGetOrderItemsByAggregateID)
}

func NewGetOrderItemsByAggregateID(sb *supabase.Client) GetOrderItemsByAggregateID {
	return func(ctx context.Context, menuID string, aggregateID int64) ([]OrderItemByAggregateRow, error) {
		data, _, err := sb.From("order_items_projection").
			Select("item_name,quantity,unit,total_price,station", "", false).
			Eq("menu_id", menuID).
			Eq("aggregate_id", strconv.FormatInt(aggregateID, 10)).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("querying order_items_projection: %w", err)
		}
		var rows []OrderItemByAggregateRow
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("unmarshaling order_items_projection: %w", err)
		}
		return rows, nil
	}
}
