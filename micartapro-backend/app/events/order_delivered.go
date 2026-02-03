package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// OrderDeliveredRequest representa el evento cuando Expediter marca que entregó
type OrderDeliveredRequest struct {
	AggregateID int64  `json:"aggregateId"`
	UpdatedAt   string `json:"updatedAt"`
}

func (o OrderDeliveredRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("order.delivered")
	event.SetType(EventOrderDelivered)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, o)
	event.SetTime(time.Now())
	return event
}
