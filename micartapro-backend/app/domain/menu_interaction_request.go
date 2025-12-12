package domain

import (
	"time"

	"github.com/alpkeskin/gotoon"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type ChatMessage struct {
	Role    string
	Content string
}

type MenuPreferences struct {
	Language string
}

type MenuInteractionRequest struct {
	MenuID          string `json:"menuId"`
	Message         string `json:"message"`
	History         []ChatMessage
	MenuPreferences MenuPreferences
	JsonMenu        MenuCreateRequest
}

func (m MenuInteractionRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("menu.interaction.request")
	event.SetType(EventMenuInteractionRequested)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, m)
	event.SetTime(time.Now())
	return event
}

func (m MenuInteractionRequest) MenuToon() string {
	encoded, _ := gotoon.Encode(m.JsonMenu)
	return encoded
}
