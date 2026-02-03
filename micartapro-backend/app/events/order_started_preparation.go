package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// OrderStartedPreparationRequest representa el evento cuando Cocina/Barra marca que está preparando
type OrderStartedPreparationRequest struct {
	AggregateID int64    `json:"aggregateId"`
	ItemKeys    []string `json:"itemKeys"` // Si está vacío, aplica a todos los items de la estación
	Station     string   `json:"station"` // KITCHEN | BAR
	UpdatedAt   string   `json:"updatedAt"`
}

func (o OrderStartedPreparationRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("order.started.preparation")
	event.SetType(EventOrderStartedPreparation)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, o)
	event.SetTime(time.Now())
	return event
}
