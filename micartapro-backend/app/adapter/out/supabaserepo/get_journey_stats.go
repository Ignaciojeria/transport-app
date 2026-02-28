package supabaserepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
)

// GetJourneyStats obtiene las estadísticas de productos para una jornada.
// Totaliza por estado: entregadas (DISPATCHED/DELIVERED), pendientes, canceladas.
// La venta se concreta solo cuando el producto está entregado.
type GetJourneyStats func(ctx context.Context, journeyID string) (*JourneyStatsResult, error)

// JourneyStatsResult es el resultado de GetJourneyStats.
type JourneyStatsResult struct {
	TotalRevenue    float64
	TotalCost       float64
	TotalOrders     int
	ItemsOrdered    int
	Products        []ProductStat
	RevenueByStatus RevenueByStatus
	OrdersByStatus  OrdersByStatus
}

// RevenueByStatus ventas totalizadas por estado.
type RevenueByStatus struct {
	Delivered  float64 `json:"delivered"`  // DELIVERED (retiro)
	Dispatched float64 `json:"dispatched"` // DISPATCHED (envío)
	Pending    float64 `json:"pending"`
	Cancelled  float64 `json:"cancelled"`
}

// OrdersByStatus órdenes totalizadas por estado.
type OrdersByStatus struct {
	Delivered  int `json:"delivered"`
	Dispatched int `json:"dispatched"`
	Pending    int `json:"pending"`
	Cancelled  int `json:"cancelled"`
}

// ProductStat es un producto con sus estadísticas (solo entregados).
type ProductStat struct {
	ProductName          string  `json:"productName"`
	QuantitySold         int     `json:"quantitySold"`
	TotalRevenue         float64 `json:"totalRevenue"`
	TotalCost            float64 `json:"totalCost"`
	Percentage           float64 `json:"percentage"`
	PercentageByQuantity float64 `json:"percentageByQuantity"`
}

// rpcStatsResponse es la respuesta cruda del RPC.
type rpcStatsResponse struct {
	RevenueByStatus struct {
		Delivered  float64 `json:"delivered"`
		Dispatched float64 `json:"dispatched"`
		Pending    float64 `json:"pending"`
		Cancelled  float64 `json:"cancelled"`
	} `json:"revenueByStatus"`
	OrdersByStatus struct {
		Delivered  int `json:"delivered"`
		Dispatched int `json:"dispatched"`
		Pending    int `json:"pending"`
		Cancelled  int `json:"cancelled"`
	} `json:"ordersByStatus"`
	Products []struct {
		ProductName  string  `json:"productName"`
		QuantitySold int     `json:"quantitySold"`
		TotalRevenue float64 `json:"totalRevenue"`
		TotalCost    float64 `json:"totalCost"`
	} `json:"products"`
	TotalRevenue float64 `json:"totalRevenue"`
	TotalCost    float64 `json:"totalCost"`
	TotalOrders  int     `json:"totalOrders"`
	ItemsOrdered int     `json:"itemsOrdered"`
}

func init() {
	ioc.Register(NewGetJourneyStats)
}

func NewGetJourneyStats(conf configuration.Conf) GetJourneyStats {
	return func(ctx context.Context, journeyID string) (*JourneyStatsResult, error) {
		rpcParams := map[string]interface{}{"p_journey_id": journeyID}
		requestBody, err := json.Marshal(rpcParams)
		if err != nil {
			return nil, fmt.Errorf("marshaling rpc params: %w", err)
		}

		rpcURL := fmt.Sprintf("%s/rest/v1/rpc/get_journey_stats_by_status", conf.SUPABASE_PROJECT_URL)
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
			return nil, fmt.Errorf("calling get_journey_stats_by_status: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("RPC failed status %d: %s", resp.StatusCode, string(body))
		}

		// PostgREST devuelve jsonb: puede ser objeto directo {...} o array [result]
		var rpc rpcStatsResponse
		if err := json.Unmarshal(body, &rpc); err != nil {
			var raw []json.RawMessage
			if err2 := json.Unmarshal(body, &raw); err2 != nil {
				return nil, fmt.Errorf("unmarshaling rpc response: %w", err)
			}
			if len(raw) == 0 {
				return &JourneyStatsResult{
					RevenueByStatus: RevenueByStatus{},
					OrdersByStatus:  OrdersByStatus{},
					Products:        []ProductStat{},
				}, nil
			}
			if err := json.Unmarshal(raw[0], &rpc); err != nil {
				return nil, fmt.Errorf("unmarshaling stats: %w", err)
			}
		}

		products := make([]ProductStat, 0, len(rpc.Products))
		var totalQuantity int
		for _, p := range rpc.Products {
			totalQuantity += p.QuantitySold
			products = append(products, ProductStat{
				ProductName:  p.ProductName,
				QuantitySold: p.QuantitySold,
				TotalRevenue: p.TotalRevenue,
				TotalCost:    p.TotalCost,
			})
		}

		for i := range products {
			if rpc.TotalRevenue > 0 {
				products[i].Percentage = (products[i].TotalRevenue / rpc.TotalRevenue) * 100
			}
			if totalQuantity > 0 {
				products[i].PercentageByQuantity = (float64(products[i].QuantitySold) / float64(totalQuantity)) * 100
			}
		}

		return &JourneyStatsResult{
			TotalRevenue: rpc.TotalRevenue,
			TotalCost:    rpc.TotalCost,
			TotalOrders:  rpc.TotalOrders,
			ItemsOrdered: rpc.ItemsOrdered,
			Products:     products,
			RevenueByStatus: RevenueByStatus{
				Delivered:  rpc.RevenueByStatus.Delivered,
				Dispatched: rpc.RevenueByStatus.Dispatched,
				Pending:    rpc.RevenueByStatus.Pending,
				Cancelled:  rpc.RevenueByStatus.Cancelled,
			},
			OrdersByStatus: OrdersByStatus{
				Delivered:  rpc.OrdersByStatus.Delivered,
				Dispatched: rpc.OrdersByStatus.Dispatched,
				Pending:    rpc.OrdersByStatus.Pending,
				Cancelled:  rpc.OrdersByStatus.Cancelled,
			},
		}, nil
	}
}
