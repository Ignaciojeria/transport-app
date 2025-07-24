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

type UpsertNodeType func(context.Context, domain.NodeType, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertNodeType, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertNodeType(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertNodeType {
	return func(ctx context.Context, nt domain.NodeType, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.NodeType

			err := tx.WithContext(ctx).
				Table("node_types").
				Where("document_id = ?", nt.DocID(ctx)).
				First(&existing).Error

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
			
			newRecord := mapper.MapNodeType(ctx, nt)

			err = tx.Omit("Tenant").Create(&newRecord).Error
			if err != nil {
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
