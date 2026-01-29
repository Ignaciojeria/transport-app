package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// ImageGenerationRequestEvent representa un evento para generar una imagen individual
type ImageGenerationRequestEvent struct {
	MenuID      string `json:"menuId"`
	MenuItemID  string `json:"menuItemId,omitempty"` // Vacío para cover image
	Prompt      string `json:"prompt"`
	AspectRatio string `json:"aspectRatio"`
	ImageCount  int    `json:"imageCount"`
	UploadURL   string `json:"uploadUrl"`   // URL pre-firmada para subir
	PublicURL   string `json:"publicUrl"`   // URL pública donde se guardará la imagen
	ImageType   string `json:"imageType"`   // "cover", "footer", o "item"
}

func (e ImageGenerationRequestEvent) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("image.generation.request")
	event.SetType(EventImageGenerationRequested)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, e)
	event.SetTime(time.Now())
	return event
}

// ImageEditionRequestEvent representa un evento para editar una imagen individual
type ImageEditionRequestEvent struct {
	MenuID           string `json:"menuId"`
	MenuItemID       string `json:"menuItemId,omitempty"` // Vacío para cover image
	Prompt           string `json:"prompt"`
	ReferenceImageUrl string `json:"referenceImageUrl"`
	AspectRatio      string `json:"aspectRatio"`
	ImageCount       int    `json:"imageCount"`
	UploadURL        string `json:"uploadUrl"` // URL pre-firmada para subir
	PublicURL        string `json:"publicUrl"` // URL pública donde se guardará la imagen
	ImageType        string `json:"imageType"` // "cover", "footer", o "item"
}

func (e ImageEditionRequestEvent) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("image.edition.request")
	event.SetType(EventImageEditionRequested)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, e)
	event.SetTime(time.Now())
	return event
}
