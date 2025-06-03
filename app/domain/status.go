package domain

import (
	"time"
)

const (
	StatusAvailable = "available"
	StatusPending   = "pending"
	StatusScanned   = "scanned"
	StatusPicked    = "picked"
	StatusPlanned   = "planned"
	StatusInTransit = "in_transit"
	StatusCancelled = "cancelled"
	StatusFinished  = "finished"
)

type Status struct {
	ID        int64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (oe Status) DocID() DocumentID {
	return HashInputs(oe.Status)
}
