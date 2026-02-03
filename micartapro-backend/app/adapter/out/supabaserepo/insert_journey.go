package supabaserepo

import (
	"context"
	"fmt"
	"time"

	"micartapro/app/shared/infrastructure/supabasecli"
	"micartapro/app/usecase/journey"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

// InsertJourney crea una nueva jornada OPEN para el menú. No comprueba si ya existe una abierta (eso lo hace el handler).
type InsertJourney func(ctx context.Context, menuID string, openedBy journey.OpenedBy, reason *string) (*journey.Journey, error)

func init() {
	ioc.Registry(NewInsertJourney, supabasecli.NewSupabaseClient)
}

func NewInsertJourney(supabase *supabase.Client) InsertJourney {
	return func(ctx context.Context, menuID string, openedBy journey.OpenedBy, reason *string) (*journey.Journey, error) {
		now := time.Now().UTC()
		id := uuid.New().String()

		record := map[string]interface{}{
			"id":         id,
			"menu_id":    menuID,
			"status":     string(journey.StatusOpen),
			"opened_at":  now.Format(time.RFC3339),
			"opened_by":  string(openedBy),
			"updated_at": now.Format(time.RFC3339),
		}
		if reason != nil && *reason != "" {
			record["opened_reason"] = *reason
		}

		_, _, err := supabase.From("journeys").
			Insert(record, false, "", "", "").
			Execute()
		if err != nil {
			return nil, fmt.Errorf("inserting journey: %w", err)
		}

		j := &journey.Journey{
			ID:           id,
			MenuID:       menuID,
			Status:       journey.StatusOpen,
			OpenedAt:     now,
			OpenedBy:     openedBy,
			OpenedReason: reason,
		}
		return j, nil
	}
}
