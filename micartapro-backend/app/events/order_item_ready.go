package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// OrderItemReadyRequest representa el evento cuando Cocina/Barra marca que terminó
type OrderItemReadyRequest struct {
	AggregateID int64    `json:"aggregateId"`
	ItemKeys    []string `json:"itemKeys"` // Items específicos que están listos
	Station     string   `json:"station"` // KITCHEN | BAR
	UpdatedAt   string   `json:"updatedAt"`
}

func (o OrderItemReadyRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("order.item.ready")
	event.SetType(EventOrderItemReady)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, o)
	event.SetTime(time.Now())
	return event
}
