package supabaserepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

// AssignOrdersResult representa el resultado de asignar una orden a una jornada.
type AssignOrdersResult struct {
	AggregateID int64 `json:"aggregate_id"`
	OrderNumber *int  `json:"order_number"`
	Assigned    bool  `json:"assigned"`
}

// AssignOrdersToJourney asigna órdenes pendientes a una jornada activa.
type AssignOrdersToJourney func(ctx context.Context, menuID, journeyID string, aggregateIDs []int64) ([]AssignOrdersResult, error)

func init() {
	ioc.Registry(NewAssignOrdersToJourney, supabasecli.NewSupabaseClient, configuration.NewConf)
}

func NewAssignOrdersToJourney(supabase *supabase.Client, conf configuration.Conf) AssignOrdersToJourney {
	return func(ctx context.Context, menuID, journeyID string, aggregateIDs []int64) ([]AssignOrdersResult, error) {
		rpcParams := map[string]interface{}{
			"p_menu_id":       menuID,
			"p_journey_id":    journeyID,
			"p_aggregate_ids": aggregateIDs,
		}

		rpcURL := fmt.Sprintf("%s/rest/v1/rpc/assign_orders_to_journey", conf.SUPABASE_PROJECT_URL)
		requestBody, err := json.Marshal(rpcParams)
		if err != nil {
			return nil, fmt.Errorf("marshaling rpc params: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Authorization", "Bearer "+conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Prefer", "return=representation")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("calling assign_orders_to_journey RPC: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("RPC call failed with status %d: %s", resp.StatusCode, string(body))
		}

		var results []AssignOrdersResult
		if err := json.Unmarshal(body, &results); err != nil {
			return nil, fmt.Errorf("unmarshaling assign results: %w", err)
		}
		return results, nil
	}
}
