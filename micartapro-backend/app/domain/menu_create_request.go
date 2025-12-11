package domain

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Side struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type MenuItem struct {
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Sides       []Side  `json:"sides,omitempty"`
	Price       float64 `json:"price,omitempty"`
}

type MenuCategory struct {
	Title string     `json:"title"`
	Items []MenuItem `json:"items"`
}

type BusinessInfo struct {
	BusinessName  string   `json:"businessName"`
	Whatsapp      string   `json:"whatsapp"`
	BusinessHours []string `json:"businessHours"`
}

type MenuCreateRequest struct {
	ID           string         `json:"id"`
	BusinessInfo BusinessInfo   `json:"businessInfo"`
	Menu         []MenuCategory `json:"menu"`
}

func (c MenuCreateRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("menu.create.request")
	event.SetType(EventMenuCreateRequested)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
