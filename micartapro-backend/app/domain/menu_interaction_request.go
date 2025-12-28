package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alpkeskin/gotoon"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

var (
	ErrInvalidMenuInteractionRequest = errors.New("invalid menu interaction request")
)

type ChatMessage struct {
	Role    string
	Content string
}

type MenuPreferences struct {
	Language string
}

type MenuInteractionRequest struct {
	Message         string `json:"message" example:"Actualiza el precio del Pollo con papas fritas a $7.500"`
	MenuID          string `json:"menuId" example:"01890a5d-ac96-7748-b800-303931383939"`
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

func (m MenuInteractionRequest) Validate() error {
	var problems []string

	parsedUUID, err := uuid.Parse(m.MenuID)
	if err != nil {
		problems = append(problems, "menuId must be a valid UUIDv7")
	} else if parsedUUID.Version() != 7 {
		problems = append(problems, "menuId must be a UUIDv7")
	}

	if len(problems) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidMenuInteractionRequest, strings.Join(problems, ", "))
	}

	return nil
}
