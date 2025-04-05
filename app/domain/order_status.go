package domain

import "time"

type OrderStatus struct {
	ID        int64
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

/*
func (os OrderStatus) UpdateIfChanged(newOrderStatus OrderStatus) OrderStatus {
	// Actualizar ID si es diferente de 0
	if newOrderStatus.ID != 0 {
		os.ID = newOrderStatus.ID
	}

	// Actualizar Status si no está vacío
	if newOrderStatus.Status != "" {
		os.Status = newOrderStatus.Status
	}

	// Actualizar CreatedAt si no es zero value
	if !newOrderStatus.CreatedAt.IsZero() {
		os.CreatedAt = newOrderStatus.CreatedAt
	}

	return os
}*/
