package supabaserepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"micartapro/app/usecase/journey"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// journeyRow representa una fila de la tabla journeys (para unmarshal desde Supabase).
type journeyRow struct {
	ID             string          `json:"id"`
	MenuID         string          `json:"menu_id"`
	Status         string          `json:"status"`
	OpenedAt       string          `json:"opened_at"`
	ClosedAt       *string         `json:"closed_at"`
	OpenedBy       string          `json:"opened_by"`
	OpenedReason   *string         `json:"opened_reason"`
	TotalsSnapshot json.RawMessage `json:"totals_snapshot"`
	ReportPDFURL   *string         `json:"report_pdf_url"`
	ReportXLSXURL  *string         `json:"report_xlsx_url"`
}

// GetActiveJourneyByMenuID devuelve la jornada con status OPEN para el menu_id, o nil si no hay.
type GetActiveJourneyByMenuID func(ctx context.Context, menuID string) (*journey.Journey, error)

func init() {
	ioc.Register(NewGetActiveJourneyByMenuID)
}

func NewGetActiveJourneyByMenuID(supabase *supabase.Client) GetActiveJourneyByMenuID {
	return func(ctx context.Context, menuID string) (*journey.Journey, error) {
		var rows []journeyRow
		data, _, err := supabase.From("journeys").
			Select("id,menu_id,status,opened_at,closed_at,opened_by,opened_reason,totals_snapshot,report_pdf_url,report_xlsx_url", "", false).
			Eq("menu_id", menuID).
			Eq("status", "OPEN").
			Limit(1, "").
			Execute()
		if err != nil {
			if err.Error() == "PGRST116" || err.Error() == "no rows in result set" {
				return nil, nil
			}
			return nil, fmt.Errorf("querying journeys: %w", err)
		}
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("unmarshaling journeys: %w", err)
		}
		if len(rows) == 0 {
			return nil, nil
		}
		r := rows[0]
		j, err := mapRowToJourney(r)
		if err != nil {
			return nil, err
		}
		return j, nil
	}
}

func mapRowToJourney(r journeyRow) (*journey.Journey, error) {
	openedAt, err := time.Parse(time.RFC3339, r.OpenedAt)
	if err != nil {
		return nil, fmt.Errorf("opened_at: %w", err)
	}
	var closedAt *time.Time
	if r.ClosedAt != nil && *r.ClosedAt != "" {
		t, err := time.Parse(time.RFC3339, *r.ClosedAt)
		if err != nil {
			return nil, fmt.Errorf("closed_at: %w", err)
		}
		closedAt = &t
	}
	j := &journey.Journey{
		ID:            r.ID,
		MenuID:        r.MenuID,
		Status:        journey.Status(r.Status),
		OpenedAt:      openedAt,
		ClosedAt:      closedAt,
		OpenedBy:      journey.OpenedBy(r.OpenedBy),
		OpenedReason:  r.OpenedReason,
		ReportPDFURL:  r.ReportPDFURL,
		ReportXLSXURL: r.ReportXLSXURL,
	}
	if len(r.TotalsSnapshot) > 0 {
		var snap journey.TotalsSnapshot
		if err := json.Unmarshal(r.TotalsSnapshot, &snap); err != nil {
			return nil, fmt.Errorf("totals_snapshot: %w", err)
		}
		j.TotalsSnapshot = &snap
	}
	return j, nil
}
