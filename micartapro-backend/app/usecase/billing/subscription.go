package billing

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	UserID uuid.UUID

	Provider string

	SubscriptionID string
	CustomerID     string
	ProductID      string

	Status string

	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	CancelAt           *time.Time
	CanceledAt         *time.Time

	Metadata map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
}
