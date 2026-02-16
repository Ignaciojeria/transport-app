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

// ReleaseOrdersFromJourney libera órdenes pendientes de una jornada (journey_id = NULL)
// para que puedan asignarse a la próxima jornada.
type ReleaseOrdersFromJourney func(ctx context.Context, journeyID string) error

func init() {
	ioc.Registry(NewReleaseOrdersFromJourney, supabasecli.NewSupabaseClient, configuration.NewConf)
}

func NewReleaseOrdersFromJourney(_ *supabase.Client, conf configuration.Conf) ReleaseOrdersFromJourney {
	return func(ctx context.Context, journeyID string) error {
		rpcParams := map[string]interface{}{"p_journey_id": journeyID}
		requestBody, err := json.Marshal(rpcParams)
		if err != nil {
			return fmt.Errorf("marshaling rpc params: %w", err)
		}

		rpcURL := fmt.Sprintf("%s/rest/v1/rpc/release_orders_from_journey", conf.SUPABASE_PROJECT_URL)
		req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Authorization", "Bearer "+conf.SUPABASE_BACKEND_API_KEY)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("calling release_orders_from_journey RPC: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response: %w", err)
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("release_orders_from_journey RPC failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil
	}
}
