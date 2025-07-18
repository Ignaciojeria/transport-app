package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type GetLastFSMTransitionByIdempotencyKey func(ctx context.Context, workflow string, entityID string) (domain.FSMState, error)

func init() {
	ioc.Registry(NewGetLastFSMTransitionByIdempotencyKey, database.NewConnectionFactory)
}

func NewGetLastFSMTransitionByIdempotencyKey(conn database.ConnectionFactory) GetLastFSMTransitionByIdempotencyKey {
	return func(ctx context.Context, idempotencyKey string, workflow string) (domain.FSMState, error) {
		var lastTransition table.FSMStateHistory

		err := conn.DB.WithContext(ctx).
			Where("workflow = ? AND idempotency_key = ?", workflow, idempotencyKey).
			Order("created_at DESC").
			First(&lastTransition).Error

		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			return domain.FSMState{}, err
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.FSMState{}, err
		}

		fsmState := domain.FSMState{
			TraceID:   lastTransition.TraceID,
			Workflow:  lastTransition.Workflow,
			State:     lastTransition.State,
			CreatedAt: lastTransition.CreatedAt,
		}

		return fsmState, nil
	}
}
