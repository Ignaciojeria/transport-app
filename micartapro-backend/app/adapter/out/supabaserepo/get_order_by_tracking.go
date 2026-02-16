package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

// ErrOrderNotFound se devuelve cuando no existe orden para el tracking_id.
var ErrOrderNotFound = errors.New("order not found for tracking_id")

// OrderTrackingItem es un ítem de la orden para la vista de tracking.
type OrderTrackingItem struct {
	ItemName   string  `json:"itemName"`
	Quantity   int     `json:"quantity"`
	Unit       string  `json:"unit"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"totalPrice"`
	Station    *string `json:"station,omitempty"`
}

// OrderByTrackingResult es el resultado de consultar una orden por tracking_id.
type OrderByTrackingResult struct {
	TrackingID   string              `json:"trackingId"`
	AggregateID  int64               `json:"aggregateId"`
	OrderNumber  int                 `json:"orderNumber"`
	MenuID       string              `json:"menuId"`
	Fulfillment  string              `json:"fulfillment"`
	JourneyID    *string             `json:"journeyId,omitempty"` // NULL si no hay jornada activa (negocio cerrado)
	Items        []OrderTrackingItem `json:"items"`
	RequestedAt  *string             `json:"requestedAt,omitempty"`
	CreatedAt    string              `json:"createdAt"`
}

type orderTrackingRow struct {
	AggregateID int64 `json:"aggregate_id"`
}

type orderItemProjectionRow struct {
	OrderNumber   int     `json:"order_number"`
	MenuID        string  `json:"menu_id"`
	JourneyID     *string `json:"journey_id"`
	ItemName      string  `json:"item_name"`
	Quantity      int     `json:"quantity"`
	Unit          string  `json:"unit"`
	Status        string  `json:"status"`
	TotalPrice    float64 `json:"total_price"`
	Station       *string `json:"station"`
	Fulfillment   string  `json:"fulfillment"`
	RequestedTime *string `json:"requested_time"`
	CreatedAt     string  `json:"created_at"`
}

// GetOrderByTrackingID obtiene el estado de una orden por su tracking_id (acceso público).
type GetOrderByTrackingID func(ctx context.Context, trackingID string) (*OrderByTrackingResult, error)

func init() {
	ioc.Registry(NewGetOrderByTrackingID, supabasecli.NewSupabaseClient)
}

func NewGetOrderByTrackingID(sb *supabase.Client) GetOrderByTrackingID {
	return func(ctx context.Context, trackingID string) (*OrderByTrackingResult, error) {
		// 1) Obtener aggregate_id desde order_tracking
		trackingData, _, err := sb.From("order_tracking").
			Select("aggregate_id", "", false).
			Eq("tracking_id", trackingID).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("querying order_tracking: %w", err)
		}
		var trackingRows []orderTrackingRow
		if err := json.Unmarshal(trackingData, &trackingRows); err != nil {
			return nil, fmt.Errorf("unmarshaling order_tracking: %w", err)
		}
		if len(trackingRows) == 0 {
			return nil, ErrOrderNotFound
		}
		aggregateID := trackingRows[0].AggregateID

		// 2) Obtener ítems desde order_items_projection
		data, _, err := sb.From("order_items_projection").
			Select("order_number,menu_id,journey_id,item_name,quantity,unit,status,total_price,station,fulfillment,requested_time,created_at", "", false).
			Eq("aggregate_id", strconv.FormatInt(aggregateID, 10)).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("querying order_items_projection: %w", err)
		}
		var rows []orderItemProjectionRow
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("unmarshaling order_items_projection: %w", err)
		}
		if len(rows) == 0 {
			return nil, ErrOrderNotFound
		}

		// Construir resultado
		items := make([]OrderTrackingItem, 0, len(rows))
		var requestedAt *string
		var createdAt string
		for i, r := range rows {
			items = append(items, OrderTrackingItem{
				ItemName:   r.ItemName,
				Quantity:   r.Quantity,
				Unit:       r.Unit,
				Status:     r.Status,
				TotalPrice: r.TotalPrice,
				Station:    r.Station,
			})
			if i == 0 {
				requestedAt = r.RequestedTime
				createdAt = r.CreatedAt
			}
		}

		return &OrderByTrackingResult{
			TrackingID:  trackingID,
			AggregateID: aggregateID,
			OrderNumber: rows[0].OrderNumber,
			MenuID:      rows[0].MenuID,
			Fulfillment: rows[0].Fulfillment,
			JourneyID:   rows[0].JourneyID,
			Items:       items,
			RequestedAt: requestedAt,
			CreatedAt:   createdAt,
		}, nil
	}
}
