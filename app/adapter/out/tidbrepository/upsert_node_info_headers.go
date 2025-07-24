package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertNodeInfoHeaders func(context.Context, domain.Headers, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertNodeInfoHeaders, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertNodeInfoHeaders(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertNodeInfoHeaders {
	return func(ctx context.Context, h domain.Headers, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			nodeInfoHeadersTbl := mapper.MapNodeInfoHeaders(ctx, h)

			err := tx.WithContext(ctx).
				Table("node_info_headers").
				Where("document_id = ?", h.DocID(ctx)).
				First(&nodeInfoHeadersTbl).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if err == nil {
				// Ya existe → solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			if err := tx.Omit("Tenant").Create(&nodeInfoHeadersTbl).Error; err != nil {
				return err
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 && saveFSMTransition != nil {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
