package tidbrepository

import (
	"context"
	"log"
	"sync"
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

type LoadStatuses func() Statuses

func NewLoadStatuses(conn database.ConnectionFactory) LoadStatuses {
	var once sync.Once
	statuses := make(Statuses)
	var records = []table.Status{
		{ID: 1, Status: domain.StatusAvailable},
		{ID: 2, Status: domain.StatusScanned},
		{ID: 3, Status: domain.StatusPicked},
		{ID: 4, Status: domain.StatusPlanned},
		{ID: 5, Status: domain.StatusInTransit},
		{ID: 6, Status: domain.StatusCancelled},
		{ID: 7, Status: domain.StatusFinished},
	}
	for _, record := range records {
		statuses[record.Status] = domain.Status{
			ID:     record.ID,
			Status: record.Status,
		}
	}
	return func() Statuses {
		once.Do(func() {
			if err := conn.WithContext(context.Background()).Save(&records).Error; err != nil {
				log.Fatalf("failed to upsert order statuses: %s", err)
			}
		})
		return statuses
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
