package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// JourneyListItem representa una jornada en el listado (para reportes).
type JourneyListItem struct {
	ID            string
	MenuID        string
	Status        string
	OpenedAt      time.Time
	ClosedAt      *time.Time
	ReportPDFURL  *string
	ReportXLSXURL *string
}

// GetJourneysByMenuID lista jornadas del menú, ordenadas por opened_at desc (más recientes primero).
type GetJourneysByMenuID func(ctx context.Context, menuID string, limit int) ([]JourneyListItem, error)

func init() {
	ioc.Register(NewGetJourneysByMenuID)
}

func NewGetJourneysByMenuID(supabase *supabase.Client) GetJourneysByMenuID {
	return func(ctx context.Context, menuID string, limit int) ([]JourneyListItem, error) {
		if limit <= 0 {
			limit = 50
		}
		data, _, err := supabase.From("journeys").
			Select("id,menu_id,status,opened_at,closed_at,report_pdf_url,report_xlsx_url", "", false).
			Eq("menu_id", menuID).
			Limit(limit, "").
			Execute()
		if err != nil {
			return nil, fmt.Errorf("querying journeys: %w", err)
		}
		var rows []journeyRow
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("unmarshaling journeys: %w", err)
		}
		out := make([]JourneyListItem, 0, len(rows))
		for _, r := range rows {
			openedAt, err := time.Parse(time.RFC3339, r.OpenedAt)
			if err != nil {
				continue
			}
			var closedAt *time.Time
			if r.ClosedAt != nil && *r.ClosedAt != "" {
				t, err := time.Parse(time.RFC3339, *r.ClosedAt)
				if err == nil {
					closedAt = &t
				}
			}
			out = append(out, JourneyListItem{
				ID:            r.ID,
				MenuID:        r.MenuID,
				Status:        r.Status,
				OpenedAt:      openedAt,
				ClosedAt:      closedAt,
				ReportPDFURL:  r.ReportPDFURL,
				ReportXLSXURL: r.ReportXLSXURL,
			})
		}
		sort.Slice(out, func(i, j int) bool { return out[i].OpenedAt.After(out[j].OpenedAt) })
		return out, nil
	}
}
