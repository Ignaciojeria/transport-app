package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alpkeskin/gotoon"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/oklog/ulid/v2"
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
	MessageID       string `json:"messageId" example:"ULID()=01G65Z755AFWAKHE12NY0CQ9FH"`
	Message         string `json:"message" example:"Actualiza el precio del Pollo con papas fritas a $7.500"`
	MenuID          string `json:"menuId" example:"ULID()=01G65Z755AFWAKHE12NY0CQ9FH"`
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

	if _, err := ulid.ParseStrict(m.MessageID); err != nil {
		problems = append(problems, "messageId must be a valid ULID")
	}

	if _, err := ulid.ParseStrict(m.MenuID); err != nil {
		problems = append(problems, "menuId must be a valid ULID")
	}

	if len(problems) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidMenuInteractionRequest, strings.Join(problems, ", "))
	}

	return nil
}
