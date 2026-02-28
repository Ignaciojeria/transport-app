package supabaserepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// PendingOrderItemRow es una fila de order_items_projection para órdenes pendientes.
type PendingOrderItemRow struct {
	AggregateID   int64   `json:"aggregate_id"`
	OrderNumber   *int    `json:"order_number"`
	MenuID        string  `json:"menu_id"`
	ItemName      string  `json:"item_name"`
	Quantity      int     `json:"quantity"`
	Unit          string  `json:"unit"`
	TotalPrice    float64 `json:"total_price"`
	Status        string  `json:"status"`
	RequestedTime *string `json:"requested_time"`
	CreatedAt     string  `json:"created_at"`
}

// PendingOrderItem es un ítem dentro de una orden pendiente (agrupado por agregado).
type PendingOrderItem struct {
	ItemName   string  `json:"item_name"`
	Quantity   int     `json:"quantity"`
	Unit       string  `json:"unit"`
	TotalPrice float64 `json:"total_price"`
}

// OrderTrackingRow es una fila de order_tracking.
type OrderTrackingRow struct {
	AggregateID int64  `json:"aggregate_id"`
	TrackingID  string `json:"tracking_id"`
}

// PendingOrder representa una orden pendiente (sin jornada) agregada por aggregate_id.
type PendingOrder struct {
	AggregateID  int64              `json:"aggregate_id"`
	TrackingID   string             `json:"tracking_id"`
	CreatedAt    string             `json:"created_at"`
	ScheduledFor *string            `json:"scheduled_for,omitempty"`
	TotalAmount  int64              `json:"total_amount"` // suma de total_price
	Status       string             `json:"status"`       // "PENDING"
	Items        []PendingOrderItem `json:"items"`        // ítems agrupados por agregado
}

// GetPendingOrdersFilter filtros opcionales para órdenes pendientes.
// FromDate y ToDate son strings UTC ISO-8601 tal cual los envía el frontend (sin transformación).
type GetPendingOrdersFilter struct {
	FromDate *string // UTC ISO-8601, ej: 2026-02-15T00:00:00.000Z
	ToDate   *string // UTC ISO-8601, ej: 2026-02-15T23:59:59.999Z
}

// GetPendingOrders obtiene órdenes con journey_id IS NULL (pendientes de asignar).
type GetPendingOrders func(ctx context.Context, menuID string, filter *GetPendingOrdersFilter) ([]PendingOrder, error)

func init() {
	ioc.Register(NewGetPendingOrders)
}

func NewGetPendingOrders(sb *supabase.Client, conf configuration.Conf) GetPendingOrders {
	return func(ctx context.Context, menuID string, filter *GetPendingOrdersFilter) ([]PendingOrder, error) {
		// RPC: BETWEEN en SQL (created_at >= p_from_date AND created_at <= p_to_date)
		rpcParams := map[string]interface{}{"p_menu_id": menuID}
		if filter != nil && filter.FromDate != nil && *filter.FromDate != "" {
			rpcParams["p_from_date"] = *filter.FromDate
		}
		if filter != nil && filter.ToDate != nil && *filter.ToDate != "" {
			rpcParams["p_to_date"] = *filter.ToDate
		}

		requestBody, err := json.Marshal(rpcParams)
		if err != nil {
			return nil, fmt.Errorf("marshaling rpc params: %w", err)
		}

		rpcURL := fmt.Sprintf("%s/rest/v1/rpc/get_pending_orders", conf.SUPABASE_PROJECT_URL)
		req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Authorization", "Bearer "+conf.SUPABASE_BACKEND_API_KEY)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("calling get_pending_orders RPC: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("get_pending_orders RPC failed with status %d: %s", resp.StatusCode, string(body))
		}

		var rows []PendingOrderItemRow
		if err := json.Unmarshal(body, &rows); err != nil {
			return nil, fmt.Errorf("unmarshaling get_pending_orders response: %w", err)
		}

		// Obtener tracking_id para cada aggregate_id
		aggIDs := make(map[int64]bool)
		for _, r := range rows {
			aggIDs[r.AggregateID] = true
		}
		trackingByAgg := make(map[int64]string)
		if len(aggIDs) > 0 {
			// Query order_tracking para los aggregate_ids que tenemos
			trackingData, _, err := sb.From("order_tracking").Select("aggregate_id,tracking_id", "", false).Execute()
			if err == nil {
				var trackingRows []OrderTrackingRow
				if err := json.Unmarshal(trackingData, &trackingRows); err == nil {
					for _, t := range trackingRows {
						if aggIDs[t.AggregateID] {
							trackingByAgg[t.AggregateID] = t.TrackingID
						}
					}
				}
			}
		}

		// Agrupar por aggregate_id: sum total_price, min created_at, requested_time, y lista de ítems
		byOrder := make(map[int64]*PendingOrder)
		for _, r := range rows {
			order, ok := byOrder[r.AggregateID]
			if !ok {
				order = &PendingOrder{
					AggregateID: r.AggregateID,
					TrackingID:  trackingByAgg[r.AggregateID],
					CreatedAt:   r.CreatedAt,
					TotalAmount: 0,
					Status:      "PENDING",
					Items:       []PendingOrderItem{},
				}
				if r.RequestedTime != nil && *r.RequestedTime != "" {
					order.ScheduledFor = r.RequestedTime
				}
				byOrder[r.AggregateID] = order
			}
			order.TotalAmount += int64(r.TotalPrice)
			order.Items = append(order.Items, PendingOrderItem{
				ItemName:   r.ItemName,
				Quantity:   r.Quantity,
				Unit:       r.Unit,
				TotalPrice: r.TotalPrice,
			})
		}

		var result []PendingOrder
		for _, order := range byOrder {
			result = append(result, *order)
		}

		// Ordenar por created_at ASC
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt < result[j].CreatedAt
		})
		return result, nil
	}
}
