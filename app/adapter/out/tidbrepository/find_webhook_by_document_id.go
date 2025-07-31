package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	"github.com/cockroachdb/errors"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewFindWebhookByDocumentID,
		database.NewConnectionFactory,
	)
}

type FindWebhookByDocumentID func(
	ctx context.Context,
	documentID domain.DocumentID,
) (domain.Webhook, error)

func NewFindWebhookByDocumentID(conn database.ConnectionFactory) FindWebhookByDocumentID {
	return func(ctx context.Context, documentID domain.DocumentID) (domain.Webhook, error) {
		// Crear una variable para almacenar el webhook desde la tabla
		var tableWebhook table.Webhook

		// Intentar buscar el webhook por document_id
		err := conn.DB.WithContext(ctx).
			Where("document_id = ?", documentID).
			First(&tableWebhook).Error

		// Manejar casos de error
		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Webhook{}, errors.Wrap(ErrWebhookNotFound, "webhook with document_id not found")
			}

			// Retornar cualquier otro error
			return domain.Webhook{}, err
		}

		webhook, err := tableWebhook.Map()
		if err != nil {
			return domain.Webhook{}, err
		}
		return webhook, nil
	}
}
