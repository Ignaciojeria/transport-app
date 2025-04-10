package domain

import (
	"strconv"
	"time"
)

type OrderStatus struct {
	ID        int64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (oe OrderStatus) DocID() DocumentID {
	return DocumentID(strconv.Itoa(int(oe.ID)))
}
