package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// OrderDispatchedRequest representa el evento cuando Expediter marca que entregó
type OrderDispatchedRequest struct {
	AggregateID int64  `json:"aggregateId"`
	UpdatedAt   string `json:"updatedAt"`
}

func (o OrderDispatchedRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("order.dispatched")
	event.SetType(EventOrderDispatched)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, o)
	event.SetTime(time.Now())
	return event
}
