package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreateOrderRequest struct {
	BusinessInfo struct {
		BusinessName string `json:"businessName"`
		Whatsapp     string `json:"whatsapp"`
	} `json:"businessInfo"`
	Items []struct {
		ProductName string  `json:"productName"`
		PricingMode string  `json:"pricingMode"`
		Unit        string  `json:"unit"`
		Quantity    float64 `json:"quantity"`
		UnitPrice   int     `json:"unitPrice"`
		TotalPrice  int     `json:"totalPrice"`
	} `json:"items"`
	Totals struct {
		Subtotal    int    `json:"subtotal"`
		DeliveryFee int    `json:"deliveryFee"`
		Total       int    `json:"total"`
		Currency    string `json:"currency"`
	} `json:"totals"`
	Fulfillment struct {
		Type          string    `json:"type"`
		RequestedTime time.Time `json:"requestedTime"`
	} `json:"fulfillment"`
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
