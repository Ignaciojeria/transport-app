package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
)

type GetFSMTransitionHistory func(ctx context.Context) ([]domain.FSMState, error)

func init() {
	ioc.Registry(NewGetFSMTransitionHistory, database.NewConnectionFactory)
}

func NewGetFSMTransitionHistory(conn database.ConnectionFactory) GetFSMTransitionHistory {
	return func(ctx context.Context) ([]domain.FSMState, error) {
		var history []table.FSMStateHistory

		traceID := trace.SpanContextFromContext(ctx).TraceID().String()

		err := conn.DB.WithContext(ctx).
			Where("trace_id = ?", traceID).
			Order("created_at ASC").
			Find(&history).Error

		if err != nil {
			return nil, err
		}

		// Mapear de tabla a dominio
		domainHistory := make([]domain.FSMState, len(history))
		for i, transition := range history {
			domainHistory[i] = domain.FSMState{
				TraceID:   transition.TraceID,
				Workflow:  transition.Workflow,
				State:     transition.State,
				CreatedAt: transition.CreatedAt,
			}
		}

		return domainHistory, nil
	}
}
