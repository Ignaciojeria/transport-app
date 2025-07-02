package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewLoadStatuses,
		database.NewConnectionFactory)
}

type LoadStatuses func() error

func NewLoadStatuses(conn database.ConnectionFactory) LoadStatuses {
	return func() error {
		available := domain.Status{
			ID:     1,
			Status: domain.StatusAvailable,
		}
		scanned := domain.Status{
			ID:     2,
			Status: domain.StatusScanned,
		}
		picked := domain.Status{
			ID:     3,
			Status: domain.StatusPicked,
		}
		planned := domain.Status{
			ID:     4,
			Status: domain.StatusPlanned,
		}
		inTransit := domain.Status{
			ID:     5,
			Status: domain.StatusInTransit,
		}
		cancelled := domain.Status{
			ID:     6,
			Status: domain.StatusCancelled,
		}
		finished := domain.Status{
			ID:     7,
			Status: domain.StatusFinished,
		}
		pending := domain.Status{
			ID:     8,
			Status: domain.StatusPending,
		}
		var records = []table.Status{
			{ID: 1, Status: available.Status, DocumentID: available.DocID().String()},
			{ID: 2, Status: scanned.Status, DocumentID: scanned.DocID().String()},
			{ID: 3, Status: picked.Status, DocumentID: picked.DocID().String()},
			{ID: 4, Status: planned.Status, DocumentID: planned.DocID().String()},
			{ID: 5, Status: inTransit.Status, DocumentID: inTransit.DocID().String()},
			{ID: 6, Status: cancelled.Status, DocumentID: cancelled.DocID().String()},
			{ID: 7, Status: finished.Status, DocumentID: finished.DocID().String()},
			{ID: 8, Status: pending.Status, DocumentID: pending.DocID().String()},
		}
		if conn.Strategy == "disabled" {
			return nil
		}
		return conn.WithContext(context.Background()).Save(&records).Error
	}
}
