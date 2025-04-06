package domain

import (
	"time"
)

type OrderStatus struct {
	ID        int64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (oe OrderStatus) DocID() DocumentID {
	return Hash(GlobalOrganization, oe.Status)
}
