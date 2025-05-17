package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	"gorm.io/gorm"
)

type UpsertState func(ctx context.Context, state domain.State) error

func NewUpsertState(db database.ConnectionFactory) UpsertState {
	return func(ctx context.Context, state domain.State) error {
		var existing table.State
		err := db.WithContext(ctx).
			Table("states").
			Where("document_id = ?", state.DocID(ctx).String()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Si ya existe, no hacemos nada porque el nombre es el mismo
		if err == nil {
			return nil
		}

		// No existe → insert
		newState := mapper.MapStateTable(ctx, state)
		return db.Create(&newState).Error
	}
}
