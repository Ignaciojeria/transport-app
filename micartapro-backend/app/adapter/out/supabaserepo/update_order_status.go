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
	"github.com/supabase-community/supabase-go"
)

type UpdateOrderStatus func(ctx context.Context, aggregateID int64, newStatus string, itemKeys []string, station string, eventType string, eventPayload interface{}) error

func init() {
	ioc.Register(NewUpdateOrderStatus)
}

func NewUpdateOrderStatus(supabase *supabase.Client, conf configuration.Conf) UpdateOrderStatus {
	return func(ctx context.Context, aggregateID int64, newStatus string, itemKeys []string, station string, eventType string, eventPayload interface{}) error {
		// Convertir eventPayload a JSON
		eventPayloadBytes, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		// Preparar los parámetros para la función RPC
		rpcParams := map[string]interface{}{
			"p_aggregate_id":  aggregateID,
			"p_new_status":    newStatus,
			"p_item_keys":     itemKeys,
			"p_station":       station,
			"p_event_type":    eventType,
			"p_event_payload": json.RawMessage(eventPayloadBytes),
		}

		// Llamar a la función stored procedure usando HTTP POST directo
		rpcURL := fmt.Sprintf("%s/rest/v1/rpc/update_order_status", conf.SUPABASE_PROJECT_URL)
		requestBody, err := json.Marshal(rpcParams)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Authorization", "Bearer "+conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Prefer", "return=representation")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// PostgREST devuelve 204 No Content cuando la función retorna void
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("RPC call failed with status %d: %s", resp.StatusCode, string(body))
		}

		return nil
	}
}
