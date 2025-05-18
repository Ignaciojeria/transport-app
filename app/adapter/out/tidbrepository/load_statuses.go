package tidbrepository

import (
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
			Status: domain.StatusAvailable,
		}
		scanned := domain.Status{
			Status: domain.StatusScanned,
		}
		picked := domain.Status{
			Status: domain.StatusPicked,
		}
		planned := domain.Status{
			Status: domain.StatusPlanned,
		}
		inTransit := domain.Status{
			Status: domain.StatusInTransit,
		}
		cancelled := domain.Status{
			Status: domain.StatusCancelled,
		}
		finished := domain.Status{
			Status: domain.StatusFinished,
		}

		var records = []table.Status{
			{ID: 1, Status: available.Status, DocumentID: available.DocID().String()},
			{ID: 2, Status: scanned.Status, DocumentID: scanned.DocID().String()},
			{ID: 3, Status: picked.Status, DocumentID: picked.DocID().String()},
			{ID: 4, Status: planned.Status, DocumentID: planned.DocID().String()},
			{ID: 5, Status: inTransit.Status, DocumentID: inTransit.DocID().String()},
			{ID: 6, Status: cancelled.Status, DocumentID: cancelled.DocID().String()},
			{ID: 7, Status: finished.Status, DocumentID: finished.DocID().String()},
		}

		return conn.DB.Table("statuses").Create(&records).Error
	}
}

type Statuses map[string]domain.Status

func (m Statuses) Available() domain.Status {
	return m[domain.StatusAvailable]
}

func (m Statuses) Scanned() domain.Status {
	return m[domain.StatusScanned]
}

func (m Statuses) Picked() domain.Status {
	return m[domain.StatusPicked]
}

func (m Statuses) Planned() domain.Status {
	return m[domain.StatusPlanned]
}

func (m Statuses) InTransit() domain.Status {
	return m[domain.StatusInTransit]
}

func (m Statuses) Cancelled() domain.Status {
	return m[domain.StatusCancelled]
}

func (m Statuses) Finished() domain.Status {
	return m[domain.StatusFinished]
}
