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

const (
	statusAvailable = "available"
	statusScanned   = "scanned"
	statusPlanned   = "planned"
	statusInTransit = "in_transit"
	statusCancelled = "cancelled"
	statusFinished  = "finished"
)

type LoadOrderStatuses func() orderStatuses

func NewLoadOrderStatuses(conn database.ConnectionFactory) LoadOrderStatuses {
	var once sync.Once
	statuses := make(orderStatuses)
	var records = []table.OrderStatus{
		{ID: 1, Status: statusAvailable},
		{ID: 2, Status: statusScanned},
		{ID: 3, Status: statusPlanned},
		{ID: 4, Status: statusInTransit},
		{ID: 5, Status: statusCancelled},
		{ID: 6, Status: statusFinished},
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
	return m[statusAvailable]
}

func (m orderStatuses) Scanned() domain.OrderStatus {
	return m[statusScanned]
}

func (m orderStatuses) Planned() domain.OrderStatus {
	return m[statusPlanned]
}

func (m orderStatuses) InTransit() domain.OrderStatus {
	return m[statusInTransit]
}

func (m orderStatuses) Cancelled() domain.OrderStatus {
	return m[statusCancelled]
}

func (m orderStatuses) Finished() domain.OrderStatus {
	return m[statusFinished]
}
