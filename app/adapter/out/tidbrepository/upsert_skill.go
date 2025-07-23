package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertSkill func(context.Context, domain.Skill, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertSkill, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertSkill(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertSkill {
	return func(ctx context.Context, skill domain.Skill, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.Skill

			err := tx.WithContext(ctx).
				Table("skills").
				Where("document_id = ?", skill.DocumentID(ctx)).
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

			// No existe → insert
			newSkill := table.Skill{
				Name:       string(skill),
				DocumentID: string(skill.DocumentID(ctx)),
				TenantID:   sharedcontext.TenantIDFromContext(ctx),
			}
			
			if err := tx.Create(&newSkill).Error; err != nil {
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
