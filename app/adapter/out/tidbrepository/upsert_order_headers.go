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

type UpsertOrderHeaders func(context.Context, domain.Headers, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertOrderHeaders, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertOrderHeaders(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertOrderHeaders {
	return func(ctx context.Context, h domain.Headers, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var orderHeaders table.OrderHeaders

			err := tx.WithContext(ctx).
				Table("order_headers").
				Where("document_id = ?", h.DocID(ctx)).
				First(&orderHeaders).Error

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

			// No existe → crear y persistir FSMState si está presente
			orderHeadersTbl := mapper.MapOrderHeaders(ctx, h)

			if err := tx.Omit("Organization").Create(&orderHeadersTbl).Error; err != nil {
				return err
			}

			if len(fsmState) > 0 {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
