package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreateOrderRequest struct {
	CreatedAt    string            `json:"createdAt"`
	Items        []OrderItem       `json:"items"`
	Totals       OrderTotals       `json:"totals"`
	Fulfillment  OrderFulfillment  `json:"fulfillment"`
	BusinessInfo OrderBusinessInfo `json:"businessInfo"`
}

type OrderItem struct {
	Unit        string  `json:"unit"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   int     `json:"unitPrice"`
	TotalPrice  int     `json:"totalPrice"`
	PricingMode string  `json:"pricingMode"`
	ProductName string  `json:"productName"`
	Station     string  `json:"station,omitempty"` // KITCHEN | BAR, opcional
}

type OrderTotals struct {
	Total       int    `json:"total"`
	Currency    string `json:"currency"`
	Subtotal    int    `json:"subtotal"`
	DeliveryFee int    `json:"deliveryFee"`
}

type OrderFulfillment struct {
	Type          string       `json:"type"`
	RequestedTime string       `json:"requestedTime"`
	Address       OrderAddress `json:"address,omitempty"`
	Contact       OrderContact `json:"contact,omitempty"`
}

type OrderAddress struct {
	RawAddress      string               `json:"rawAddress"`
	Coordinates     OrderCoordinates     `json:"coordinates"`
	DeliveryDetails OrderDeliveryDetails `json:"deliveryDetails,omitempty"`
}

type OrderCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type OrderDeliveryDetails struct {
	Unit  string `json:"unit"`
	Notes string `json:"notes"`
}

type OrderContact struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type OrderBusinessInfo struct {
	Whatsapp     string `json:"whatsapp"`
	BusinessName string `json:"businessName"`
}

func (c CreateOrderRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("create.order.request")
	event.SetType(EventCreateOrderRequested)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
