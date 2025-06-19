package mapper

import (
	"time"
	"transport-app/app/domain"
)

func MapOrderStatusToDomain(status struct {
	ID        int64
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}) domain.Status {
	createdAt, err := time.Parse("2006-01-02T15:04:05Z07:00", status.CreatedAt)
	if err != nil {
		createdAt = time.Time{} // zero value como default en caso de error
	}

	return domain.Status{
		ID:        status.ID,
		Status:    status.Status,
		CreatedAt: createdAt,
	}
}
