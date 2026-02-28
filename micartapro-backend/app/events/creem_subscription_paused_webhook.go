package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemSubscriptionPausedWebhook struct {
	ID        string                        `json:"id" example:"evt_5veN2cn5N9Grz8u7w3yJuL"`
	EventType string                        `json:"eventType" example:"subscription.paused"`
	CreatedAt int64                         `json:"created_at" example:"1754041946898"`
	Object    CreemSubscriptionPausedObject `json:"object"`
}

type CreemSubscriptionPausedProduct struct {
	ID                string `json:"id" example:"prod_sYwbyE1tPbsqbLu6S0bsR"`
	Object            string `json:"object" example:"product"`
	Name              string `json:"name" example:"Prod"`
	Description       string `json:"description" example:"My Product Description"`
	ImageURL          any    `json:"image_url"`
	Price             int    `json:"price" example:"2000"`
	Currency          string `json:"currency" example:"EUR"`
	BillingType       string `json:"billing_type" example:"recurring"`
	BillingPeriod     string `json:"billing_period" example:"every-month"`
	Status            string `json:"status" example:"active"`
	TaxMode           string `json:"tax_mode" example:"exclusive"`
	TaxCategory       string `json:"tax_category" example:"saas"`
	DefaultSuccessURL string `json:"default_success_url" example:""`
	CreatedAt         string `json:"created_at" example:"2025-08-01T09:51:26.277Z"`
	UpdatedAt         string `json:"updated_at" example:"2025-08-01T09:51:26.277Z"`
	Mode              string `json:"mode" example:"test"`
}

type CreemSubscriptionPausedCustomer struct {
	ID        string `json:"id" example:"cust_4fpU8kYkQmI1XKBwU2qeME"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Test Test"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2024-11-07T23:21:11.763Z"`
	UpdatedAt string `json:"updated_at" example:"2024-11-07T23:21:11.763Z"`
	Mode      string `json:"mode" example:"test"`
}

type CreemSubscriptionPausedItems struct {
	Object    string `json:"object" example:"subscription_item"`
	ID        string `json:"id" example:"sitem_1ZIqcUuxKKDTj5WZPNsN6C"`
	ProductID string `json:"product_id" example:"prod_sYwbyE1tPbsqbLu6S0bsR"`
	PriceID   string `json:"price_id" example:"pprice_1uM3Pi1vJJ3xkhwQuZiM42"`
	Units     int    `json:"units" example:"1"`
	CreatedAt string `json:"created_at" example:"2025-08-01T09:51:50.497Z"`
	UpdatedAt string `json:"updated_at" example:"2025-08-01T09:51:50.497Z"`
	Mode      string `json:"mode" example:"test"`
}

type CreemSubscriptionPausedObject struct {
	ID                     string                          `json:"id" example:"sub_3ZT1iYMeDBpiUpRTqq4veE"`
	Object                 string                          `json:"object" example:"subscription"`
	Product                CreemSubscriptionPausedProduct  `json:"product"`
	Customer               CreemSubscriptionPausedCustomer `json:"customer"`
	Items                  []CreemSubscriptionPausedItems  `json:"items"`
	CollectionMethod       string                          `json:"collection_method" example:"charge_automatically"`
	Status                 string                          `json:"status" example:"paused"`
	CurrentPeriodStartDate string                          `json:"current_period_start_date" example:"2025-08-01T09:51:47.000Z"`
	CurrentPeriodEndDate   string                          `json:"current_period_end_date" example:"2025-09-01T09:51:47.000Z"`
	CanceledAt             any                             `json:"canceled_at" example:"null"`
	CreatedAt              string                          `json:"created_at" example:"2025-08-01T09:51:50.488Z"`
	UpdatedAt              string                          `json:"updated_at" example:"2025-08-01T09:52:26.822Z"`
	Mode                   string                          `json:"mode" example:"test"`
}

func (c CreemSubscriptionPausedWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.subscription.paused.webhook")
	event.SetType(EventCreemSubscriptionPausedWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
