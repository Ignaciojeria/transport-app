package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// SaveFSMTransition es una función que guarda transiciones de FSM
type SaveFSMTransition func(ctx context.Context, fsmState domain.FSMState, tx ...*gorm.DB) error

func init() {
	ioc.Registry(NewSaveFSMTransition, database.NewConnectionFactory)
}

func NewSaveFSMTransition(conn database.ConnectionFactory) SaveFSMTransition {
	return func(ctx context.Context, fsmState domain.FSMState, tx ...*gorm.DB) error {
		record := table.FSMStateHistory{
			TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
			IdempotencyKey: fsmState.IdempotencyKey,
			Workflow:       fsmState.Workflow,
			State:          fsmState.State,
		}
		// Determinar qué conexión usar
		var db *gorm.DB
		if len(tx) > 0 && tx[0] != nil {
			// Usar la transacción proporcionada
			db = tx[0].WithContext(ctx)
		} else {
			// Usar la conexión por defecto
			db = conn.DB.WithContext(ctx)
		}
		// Guardar el registro en la base de datos
		return db.Create(&record).Error
	}
}
