package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"

	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

// JourneyProductStatRow es una fila de journey_product_stats.
type JourneyProductStatRow struct {
	JourneyID    string   `json:"journey_id"`
	ProductName  string   `json:"product_name"`
	QuantitySold int      `json:"quantity_sold"`
	TotalRevenue *float64 `json:"total_revenue"`
}

// GetJourneyStats obtiene las estadísticas de productos para una jornada.
// Devuelve totalRevenue, totalOrders (count distinct aggregate_id en order_items_projection)
// y productos con quantity_sold, total_revenue, percentage.
type GetJourneyStats func(ctx context.Context, journeyID string) (*JourneyStatsResult, error)

// JourneyStatsResult es el resultado de GetJourneyStats.
type JourneyStatsResult struct {
	TotalRevenue float64
	TotalOrders  int
	Products     []ProductStat
}

// ProductStat es un producto con sus estadísticas.
type ProductStat struct {
	ProductName          string  `json:"productName"`
	QuantitySold         int     `json:"quantitySold"`
	TotalRevenue         float64 `json:"totalRevenue"`
	Percentage           float64 `json:"percentage"`           // % sobre revenue total
	PercentageByQuantity float64 `json:"percentageByQuantity"` // % sobre unidades totales
}

func init() {
	ioc.Registry(NewGetJourneyStats, supabasecli.NewSupabaseClient)
}

func NewGetJourneyStats(supabase *supabase.Client) GetJourneyStats {
	return func(ctx context.Context, journeyID string) (*JourneyStatsResult, error) {
		// 1. Obtener stats por producto
		data, _, err := supabase.From("journey_product_stats").
			Select("journey_id,product_name,quantity_sold,total_revenue", "", false).
			Eq("journey_id", journeyID).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("querying journey_product_stats: %w", err)
		}
		var rows []JourneyProductStatRow
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("unmarshaling journey_product_stats: %w", err)
		}

		// 2. Obtener total de órdenes (distinct aggregate_id en order_items_projection)
		ordersData, _, err := supabase.From("order_items_projection").
			Select("aggregate_id", "", false).
			Eq("journey_id", journeyID).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("querying order_items_projection for order count: %w", err)
		}
		var orderRows []struct {
			AggregateID int64 `json:"aggregate_id"`
		}
		if err := json.Unmarshal(ordersData, &orderRows); err != nil {
			return nil, fmt.Errorf("unmarshaling order count: %w", err)
		}
		seen := make(map[int64]struct{})
		for _, r := range orderRows {
			seen[r.AggregateID] = struct{}{}
		}
		totalOrders := len(seen)

		// 3. Calcular totales y porcentajes
		var totalRevenue float64
		var totalQuantity int
		products := make([]ProductStat, 0, len(rows))
		for _, r := range rows {
			rev := 0.0
			if r.TotalRevenue != nil {
				rev = *r.TotalRevenue
			}
			totalRevenue += rev
			totalQuantity += r.QuantitySold
			products = append(products, ProductStat{
				ProductName:  r.ProductName,
				QuantitySold: r.QuantitySold,
				TotalRevenue: rev,
			})
		}

		// Porcentaje por producto (revenue y cantidad)
		for i := range products {
			if totalRevenue > 0 {
				products[i].Percentage = (products[i].TotalRevenue / totalRevenue) * 100
			}
			if totalQuantity > 0 {
				products[i].PercentageByQuantity = (float64(products[i].QuantitySold) / float64(totalQuantity)) * 100
			}
		}

		return &JourneyStatsResult{
			TotalRevenue: totalRevenue,
			TotalOrders:  totalOrders,
			Products:     products,
		}, nil
	}
}
