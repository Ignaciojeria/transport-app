package supabaserepo

import (
	"context"
	"fmt"
	"time"

	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

// CloseJourney cierra la jornada: status = CLOSED, closed_at = now().
// totalsSnapshot: opcional; si no es nil, se guarda en journeys.totals_snapshot.
// Solo actualiza si id y menu_id coinciden y status es OPEN.
type CloseJourney func(ctx context.Context, menuID, journeyID string, totalsSnapshot interface{}) error

func init() {
	ioc.Registry(NewCloseJourney, supabasecli.NewSupabaseClient)
}

func NewCloseJourney(supabase *supabase.Client) CloseJourney {
	return func(ctx context.Context, menuID, journeyID string, totalsSnapshot interface{}) error {
		now := time.Now().UTC()
		record := map[string]interface{}{
			"status":     "CLOSED",
			"closed_at":  now.Format(time.RFC3339),
			"updated_at": now.Format(time.RFC3339),
		}
		if totalsSnapshot != nil {
			record["totals_snapshot"] = totalsSnapshot
		}
		_, _, err := supabase.From("journeys").
			Update(record, "", "").
			Eq("id", journeyID).
			Eq("menu_id", menuID).
			Eq("status", "OPEN").
			Execute()
		if err != nil {
			return fmt.Errorf("closing journey: %w", err)
		}
		return nil
	}
}
