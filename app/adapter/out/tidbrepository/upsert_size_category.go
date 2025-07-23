package tidbrepository

import (
	"context"
	"errors"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertSizeCategory func(context.Context, domain.SizeCategory, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertSizeCategory, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertSizeCategory(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertSizeCategory {
	return func(ctx context.Context, sc domain.SizeCategory, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.SizeCategory

			err := tx.WithContext(ctx).
				Table("size_categories").
				Where("document_id = ?", sc.DocumentID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if err == nil {
				// Ya existe → solo persistir FSMState si está presente
				if len(fsmState) > 0 {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			newRecord := mapper.MapSizeCategory(ctx, sc)

			err = tx.Create(&newRecord).Error
			if err != nil {
				return err
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
