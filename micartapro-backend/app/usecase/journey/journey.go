package journey

import "time"

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusClosed Status = "CLOSED"
)

type OpenedBy string

const (
	OpenedByUser   OpenedBy = "USER"
	OpenedBySystem OpenedBy = "SYSTEM"
)

type Journey struct {
	ID     string
	MenuID string

	Status Status

	OpenedAt time.Time
	ClosedAt *time.Time

	OpenedBy     OpenedBy
	OpenedReason *string

	TotalsSnapshot *TotalsSnapshot

	ReportPDFURL  *string
	ReportXLSXURL *string
}

type TotalsSnapshot struct {
	OrdersTotal     int   `json:"orders_total"`
	OrdersDelivered int   `json:"orders_delivered"`
	OrdersCancelled int   `json:"orders_cancelled"`
	OrdersPending   int   `json:"orders_pending"`
	PickupOrders    int   `json:"pickup_orders"`
	DeliveryOrders  int   `json:"delivery_orders"`
	GrossAmount     int64 `json:"gross_amount"` // pesos / centavos
}
