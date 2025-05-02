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
		NewLoadOrderStatuses,
		database.NewConnectionFactory)
}

type LoadOrderStatuses func() orderStatuses

func NewLoadOrderStatuses(conn database.ConnectionFactory) LoadOrderStatuses {
	var once sync.Once
	statuses := make(orderStatuses)
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
		statuses[record.Status] = domain.OrderStatus{
			ID:     record.ID,
			Status: record.Status,
		}
	}
	return func() orderStatuses {
		once.Do(func() {
			if err := conn.WithContext(context.Background()).Save(&records).Error; err != nil {
				log.Fatalf("failed to upsert order statuses: %s", err)
			}
		})
		return statuses
	}
}

type orderStatuses map[string]domain.OrderStatus

func (m orderStatuses) Available() domain.OrderStatus {
	return m[domain.StatusAvailable]
}

func (m orderStatuses) Scanned() domain.OrderStatus {
	return m[domain.StatusScanned]
}

func (m orderStatuses) Picked() domain.OrderStatus {
	return m[domain.StatusPicked]
}

func (m orderStatuses) Planned() domain.OrderStatus {
	return m[domain.StatusPlanned]
}

func (m orderStatuses) InTransit() domain.OrderStatus {
	return m[domain.StatusInTransit]
}

func (m orderStatuses) Cancelled() domain.OrderStatus {
	return m[domain.StatusCancelled]
}

func (m orderStatuses) Finished() domain.OrderStatus {
	return m[domain.StatusFinished]
}
