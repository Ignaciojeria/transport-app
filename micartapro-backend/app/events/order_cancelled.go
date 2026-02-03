package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// OrderCancelledRequest representa el evento cuando se cancela una orden
type OrderCancelledRequest struct {
	AggregateID int64  `json:"aggregateId"`
	Reason      string `json:"reason,omitempty"` // Razón opcional de la cancelación
	UpdatedAt   string `json:"updatedAt"`
}

func (o OrderCancelledRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("order.cancelled")
	event.SetType(EventOrderCancelled)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, o)
	event.SetTime(time.Now())
	return event
}
