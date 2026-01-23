package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

/*
	UNITS
*/

type UnitOfMeasure string

const (
	UnitEach        UnitOfMeasure = "EACH"
	UnitGram        UnitOfMeasure = "GRAM"
	UnitKilogram    UnitOfMeasure = "KILOGRAM"
	UnitMilliliter  UnitOfMeasure = "MILLILITER"
	UnitLiter       UnitOfMeasure = "LITER"
	UnitMeter       UnitOfMeasure = "METER"
	UnitSquareMeter UnitOfMeasure = "SQUARE_METER"
)

/*
	PRICING
*/

type PricingMode string

const (
	PricingUnit   PricingMode = "UNIT"
	PricingWeight PricingMode = "WEIGHT"
	PricingVolume PricingMode = "VOLUME"
	PricingLength PricingMode = "LENGTH"
	PricingArea   PricingMode = "AREA"
)

type Pricing struct {
	Mode         PricingMode   `json:"mode"`
	Unit         UnitOfMeasure `json:"unit"`
	PricePerUnit float64       `json:"pricePerUnit"`
	BaseUnit     float64       `json:"baseUnit"`
}

/*
	MENU MODELS
*/

type Side struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Pricing  Pricing `json:"pricing"`
	PhotoUrl string  `json:"photoUrl,omitempty"`
}

type MenuItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Sides       []Side  `json:"sides,omitempty"`
	Pricing     Pricing `json:"pricing"`
	PhotoUrl    string  `json:"photoUrl,omitempty"`
}

type MenuCategory struct {
	Title string     `json:"title"`
	Items []MenuItem `json:"items"`
}

/*
	BUSINESS
*/

type BusinessInfo struct {
	BusinessName  string   `json:"businessName"`
	Whatsapp      string   `json:"whatsapp"`
	BusinessHours []string `json:"businessHours"`
}

/*
	DELIVERY OPTIONS
*/

type DeliveryOptionType string

const (
	DeliveryOptionPickup   DeliveryOptionType = "PICKUP"
	DeliveryOptionDelivery DeliveryOptionType = "DELIVERY"
)

type TimeRequestType string

const (
	TimeRequestExact  TimeRequestType = "EXACT"
	TimeRequestWindow TimeRequestType = "WINDOW"
)

type TimeWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DeliveryOption struct {
	Type            DeliveryOptionType `json:"type"`
	RequireTime     bool               `json:"requireTime"`
	TimeRequestType TimeRequestType    `json:"timeRequestType,omitempty"`
	TimeWindows     []TimeWindow       `json:"timeWindows,omitempty"`
}

/*
	IMAGE GENERATION
*/

type ImageGenerationRequest struct {
	MenuItemID string  `json:"menuItemId"`
	Prompt     string  `json:"prompt"`
	AspectRatio string `json:"aspectRatio"`
	ImageCount  int    `json:"imageCount"`
}

/*
	EVENT
*/

type MenuCreateRequest struct {
	ID                    string                   `json:"id"`
	CoverImage            string                   `json:"coverImage"`
	FooterImage           string                   `json:"footerImage"`
	BusinessInfo          BusinessInfo             `json:"businessInfo"`
	Menu                  []MenuCategory           `json:"menu"`
	DeliveryOptions       []DeliveryOption         `json:"deliveryOptions,omitempty"`
	ImageGenerationRequests []ImageGenerationRequest `json:"imageGenerationRequests,omitempty"`
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
