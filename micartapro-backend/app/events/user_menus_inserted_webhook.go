package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type UserMenusInsertedWebhook struct {
	Type   string `json:"type" example:"insert"`
	Table  string `json:"table" example:"user_menus"`
	Schema string `json:"schema" example:"public"`
	Record struct {
		ID        int64  `json:"id" example:"1"`
		UserID    string `json:"user_id" example:"uuidv7"`
		MenuID    string `json:"menu_id" example:"uuidv7"`
		CreatedAt string `json:"created_at" example:"2021-01-01T00:00:00Z"`
	} `json:"record"`
}

func (r *UserMenusInsertedWebhook) CreatedAtToISO8601() {
	if r.Record.CreatedAt == "" {
		return
	}
	// Intentar parsear el string en diferentes formatos comunes
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, r.Record.CreatedAt)
		if err == nil {
			break
		}
	}

	if err != nil {
		// Si no se puede parsear, mantener el string original
		return
	}

	// Modificar el campo CreatedAt con el formato ISO8601
	r.Record.CreatedAt = t.UTC().Format(time.RFC3339Nano)
}

func (r *UserMenusInsertedWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("user.menus.inserted.webhook") //struct name
	event.SetType(EventUserMenusInsertedWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, r)
	event.SetTime(time.Now())
	return event
}
