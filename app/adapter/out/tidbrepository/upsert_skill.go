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

type UpsertSkill func(context.Context, domain.Skill) error

func init() {
	ioc.Registry(NewUpsertSkill, database.NewConnectionFactory)
}

func NewUpsertSkill(conn database.ConnectionFactory) UpsertSkill {
	return func(ctx context.Context, skill domain.Skill) error {
		var existing table.Skill

		err := conn.DB.WithContext(ctx).
			Table("skills").
			Where("document_id = ?", skill.DocumentID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			// Ya existe → no hacer nada
			return nil
		}

		// No existe → insert
		newSkill := table.Skill{
			Name:       string(skill),
			DocumentID: string(skill.DocumentID(ctx)),
			TenantID:   sharedcontext.TenantIDFromContext(ctx),
		}
		return conn.Create(&newSkill).Error
	}
}
