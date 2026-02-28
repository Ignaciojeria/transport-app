package supabaserepo

import (
	"context"
	"fmt"
	"time"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/supabase-community/supabase-go"
)

// UpdateJourneyReportURL actualiza report_xlsx_url en la jornada.
type UpdateJourneyReportURL func(ctx context.Context, journeyID, reportURL string) error

func init() {
	ioc.Register(NewUpdateJourneyReportURL)
}

func NewUpdateJourneyReportURL(supabase *supabase.Client) UpdateJourneyReportURL {
	return func(ctx context.Context, journeyID, reportURL string) error {
		_, _, err := supabase.From("journeys").
			Update(map[string]interface{}{
				"report_xlsx_url": reportURL,
				"updated_at":      time.Now().UTC().Format(time.RFC3339),
			}, "", "").
			Eq("id", journeyID).
			Execute()
		if err != nil {
			return fmt.Errorf("updating journey report url: %w", err)
		}
		return nil
	}
}
