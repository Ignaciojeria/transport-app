package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type GetLastFSMTransition func(ctx context.Context) (domain.FSMState, error)

func init() {
	ioc.Registry(NewGetLastFSMTransition, database.NewConnectionFactory)
}

func NewGetLastFSMTransition(conn database.ConnectionFactory) GetLastFSMTransition {
	return func(ctx context.Context) (domain.FSMState, error) {
		var lastTransition table.FSMStateHistory

		traceID := trace.SpanContextFromContext(ctx).TraceID().String()

		err := conn.DB.WithContext(ctx).
			Where("trace_id = ?", traceID).
			Order("created_at DESC").
			First(&lastTransition).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.FSMState{}, domain.ErrNoTransitionsFound
			}
			return domain.FSMState{}, err
		}

		fsmState := domain.FSMState{
			TraceID:        lastTransition.TraceID,
			IdempotencyKey: lastTransition.IdempotencyKey,
			Workflow:       lastTransition.Workflow,
			State:          lastTransition.State,
			NextInput:      lastTransition.NextInput,
			CreatedAt:      lastTransition.CreatedAt,
		}

		return fsmState, nil
	}
}
