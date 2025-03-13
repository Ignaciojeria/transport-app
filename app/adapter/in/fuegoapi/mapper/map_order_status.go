package mapper

import (
	"time"
	"transport-app/app/domain"
)

func MapOrderStatusFromDomain(status domain.OrderStatus) struct {
	ID        int64
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
} {
	return struct {
		ID        int64
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`
	}{
		ID:        status.ID,
		Status:    status.Status,
		CreatedAt: status.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func MapOrderStatusToDomain(status struct {
	ID        int64
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}) domain.OrderStatus {
	createdAt, err := time.Parse("2006-01-02T15:04:05Z07:00", status.CreatedAt)
	if err != nil {
		createdAt = time.Time{} // zero value como default en caso de error
	}

	return domain.OrderStatus{
		ID:        status.ID,
		Status:    status.Status,
		CreatedAt: createdAt,
	}
}
