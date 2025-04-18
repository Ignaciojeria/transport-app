package tidbrepository

/*
import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/gcppublisher"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewSaveEventOutBox,
		database.NewConnectionFactory,
		gcppublisher.NewApplicationEvents)
}

type SaveEventOutBox func(
	context.Context,
	domain.Outbox) (domain.Outbox, error)

func NewSaveEventOutBox(
	conn database.ConnectionFactory,
	publishOutBoxEvent gcppublisher.ApplicationEvents,
) SaveEventOutBox {
	return func(ctx context.Context, event domain.Outbox) (domain.Outbox, error) {
		// Mapear al modelo de la base de datos
		e := table.MapEventOutbox(event)
		e.Status = "pending"
		// Guardar el evento en la base de datos
		if err := conn.Save(&e).Error; err != nil {
			return domain.Outbox{}, fmt.Errorf("failed to save outbox event: %w", err)
		}

		// Manejar publicación del evento en una goroutine
		go func() {
			if pubErr := publishOutBoxEvent(context.Background(), event); pubErr != nil {
				fmt.Printf("failed to publish event %d: %v\n", e.ID, pubErr)
				return
			}
		}()

		return event, nil
	}
}
*/
