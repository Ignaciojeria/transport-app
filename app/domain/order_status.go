package domain

import (
	"context"
	"time"
)

type OrderStatus struct {
	ID        int64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (oe OrderStatus) DocID() DocumentID {
	return HashByTenant(context.Background(), oe.Status)
}
