package supabaserepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"micartapro/app/shared/infrastructure/supabasecli"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

// ErrTrackingIDConflict es devuelto cuando tracking_id ya existe (colisión).
var ErrTrackingIDConflict = errors.New("tracking_id already exists")

// InsertOrderTracking inserta una fila en la proyección order_tracking.
// aggregate_id = ID del agregado de la orden. Retorna ErrTrackingIDConflict si tracking_id ya existe.
type InsertOrderTracking func(ctx context.Context, aggregateID int64, trackingID string) error

func init() {
	ioc.Registry(NewInsertOrderTracking, supabasecli.NewSupabaseClient)
}

func NewInsertOrderTracking(supabase *supabase.Client) InsertOrderTracking {
	return func(ctx context.Context, aggregateID int64, trackingID string) error {
		record := map[string]interface{}{
			"aggregate_id": aggregateID,
			"tracking_id":  trackingID,
		}
		_, _, err := supabase.From("order_tracking").
			Insert(record, false, "", "", "").
			Execute()
		if err != nil {
			// Código PostgreSQL unique_violation = 23505
			if strings.Contains(err.Error(), "23505") || strings.Contains(strings.ToLower(err.Error()), "unique") {
				return ErrTrackingIDConflict
			}
			return fmt.Errorf("inserting order_tracking: %w", err)
		}
		return nil
	}
}
