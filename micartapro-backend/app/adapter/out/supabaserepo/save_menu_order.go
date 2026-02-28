package supabaserepo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

var ErrMenuOrderNotFound = errors.New("menu order not found")

type SaveMenuOrderResult struct {
	OrderNumber int   `json:"order_number"`
	AggregateID int64 `json:"aggregate_id"`
}

type SaveMenuOrder func(ctx context.Context, menuID string, eventPayload interface{}, eventType string) (SaveMenuOrderResult, error)

func init() {
	ioc.Register(NewSaveMenuOrder)
}

func NewSaveMenuOrder(supabase *supabase.Client, conf configuration.Conf) SaveMenuOrder {
	return func(ctx context.Context, menuID string, eventPayload interface{}, eventType string) (SaveMenuOrderResult, error) {
		// Convertir eventPayload a JSON
		eventPayloadBytes, err := json.Marshal(eventPayload)
		if err != nil {
			return SaveMenuOrderResult{}, err
		}

		// Preparar los parámetros para la función RPC
		// La función stored procedure 'create_menu_order' debe existir en la base de datos
		// Ver: sql/create_menu_order_function.sql
		rpcParams := map[string]interface{}{
			"p_menu_id":       menuID,
			"p_event_payload": json.RawMessage(eventPayloadBytes),
			"p_event_type":    eventType,
		}

		// Llamar a la función stored procedure de forma atómica usando HTTP POST directo
		// PostgREST expone las funciones RPC en /rest/v1/rpc/function_name
		rpcURL := fmt.Sprintf("%s/rest/v1/rpc/create_menu_order", conf.SUPABASE_PROJECT_URL)
		requestBody, err := json.Marshal(rpcParams)
		if err != nil {
			return SaveMenuOrderResult{}, err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return SaveMenuOrderResult{}, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Authorization", "Bearer "+conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Prefer", "return=representation")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return SaveMenuOrderResult{}, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return SaveMenuOrderResult{}, err
		}

		if resp.StatusCode != http.StatusOK {
			return SaveMenuOrderResult{}, fmt.Errorf("RPC call failed with status %d: %s", resp.StatusCode, string(body))
		}

		// La función RPC retorna un array con order_number y aggregate_id
		// Formato: [{"order_number": 1, "aggregate_id": 100023}]
		var resultArray []SaveMenuOrderResult
		if err := json.Unmarshal(body, &resultArray); err != nil {
			return SaveMenuOrderResult{}, fmt.Errorf("failed to parse RPC response: %w, body: %s", err, string(body))
		}

		if len(resultArray) == 0 {
			return SaveMenuOrderResult{}, fmt.Errorf("RPC returned empty result, body: %s", string(body))
		}

		result := resultArray[0]

		return result, nil
	}
}
