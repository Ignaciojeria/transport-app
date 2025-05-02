package domain

import (
	"context"
	"time"
)

const (
	StatusAvailable = "available"
	StatusScanned   = "scanned"
	StatusPicked    = "picked"
	StatusPlanned   = "planned"
	StatusInTransit = "in_transit"
	StatusCancelled = "cancelled"
	StatusFinished  = "finished"
)

type OrderStatus struct {
	ID        int64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (oe OrderStatus) DocID() DocumentID {
	return HashByTenant(context.Background(), oe.Status)
}
