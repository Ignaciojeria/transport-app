package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemSubscriptionTrialingWebhook struct {
	ID        string                          `json:"id" example:"evt_1NX6xlPtWIgAx7NtUegwO8"`
	EventType string                          `json:"eventType" example:"subscription.trialing"`
	CreatedAt int64                           `json:"created_at" example:"1767404212202"`
	Object    CreemSubscriptionTrialingObject `json:"object"`
}

type CreemSubscriptionTrialingProduct struct {
	ID                string `json:"id" example:"prod_amLs4z3rWJvj4ZC6xf8if"`
	Object            string `json:"object" example:"product"`
	Name              string `json:"name" example:"micartapro"`
	Description       string `json:"description" example:"Create and manage digital menus in seconds using AI. Easily update prices, items, and descriptions, and share your menu instantly via link."`
	ImageURL          string `json:"image_url" example:"https://nucn5fajkcc6sgrd.public.blob.vercel-storage.com/micartapro-cover-r2Jp8LoV97OjyqZJ5HAhtZFNdlGKuP.png"`
	Price             int    `json:"price" example:"1500"`
	Currency          string `json:"currency" example:"USD"`
	BillingType       string `json:"billing_type" example:"recurring"`
	BillingPeriod     string `json:"billing_period" example:"every-month"`
	Status            string `json:"status" example:"active"`
	TaxMode           string `json:"tax_mode" example:"exclusive"`
	TaxCategory       string `json:"tax_category" example:"saas"`
	DefaultSuccessURL any    `json:"default_success_url" example:"null"`
	CreatedAt         string `json:"created_at" example:"2025-12-30T00:47:53.004Z"`
	UpdatedAt         string `json:"updated_at" example:"2025-12-30T00:47:53.004Z"`
	Mode              string `json:"mode" example:"test"`
}

type CreemSubscriptionTrialingCustomer struct {
	ID        string `json:"id" example:"cust_2HMB2DRr4gNDtKdodcx509"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"ignaciovl.j@gmail.com"`
	Name      string `json:"name" example:"Ignacio Jeria"`
	Country   string `json:"country" example:"CL"`
	CreatedAt string `json:"created_at" example:"2025-12-30T02:12:54.604Z"`
	UpdatedAt string `json:"updated_at" example:"2025-12-30T02:12:54.604Z"`
	Mode      string `json:"mode" example:"test"`
}

type CreemSubscriptionTrialingItems struct {
	Object    string `json:"object" example:"subscription_item"`
	ID        string `json:"id" example:"sitem_a1RP5Ou3cWV0J4h0iaBcN"`
	ProductID string `json:"product_id" example:"prod_amLs4z3rWJvj4ZC6xf8if"`
	PriceID   string `json:"price_id" example:"pprice_3IYhTKGGR5fZLxjjyFqbZK"`
	Units     int    `json:"units" example:"1"`
	CreatedAt string `json:"created_at" example:"2026-01-03T01:36:52.124Z"`
	UpdatedAt string `json:"updated_at" example:"2026-01-03T01:36:52.124Z"`
	Mode      string `json:"mode" example:"test"`
}
type CreemSubscriptionTrialingMetadata struct {
	UserID string `json:"user_id" example:"763a590a-9b8e-4a91-b8ee-47f2a64d003d"`
}
type CreemSubscriptionTrialingObject struct {
	ID                     string                            `json:"id" example:"sub_6FjNokyn7Zvl9TxtogTo3q"`
	Object                 string                            `json:"object" example:"subscription"`
	Product                CreemSubscriptionTrialingProduct  `json:"product"`
	Customer               CreemSubscriptionTrialingCustomer `json:"customer"`
	Items                  []CreemSubscriptionTrialingItems  `json:"items"`
	CollectionMethod       string                            `json:"collection_method" example:"charge_automatically"`
	Status                 string                            `json:"status" example:"trialing"`
	CurrentPeriodStartDate string                            `json:"current_period_start_date" example:"2026-01-03T01:36:48.000Z"`
	CurrentPeriodEndDate   string                            `json:"current_period_end_date" example:"2026-01-17T01:36:48.000Z"`
	CanceledAt             any                               `json:"canceled_at" example:"null"`
	CreatedAt              string                            `json:"created_at" example:"2026-01-03T01:36:52.097Z"`
	UpdatedAt              string                            `json:"updated_at" example:"2026-01-03T01:36:52.097Z"`
	Metadata               CreemSubscriptionTrialingMetadata `json:"metadata"`
}

func (c CreemSubscriptionTrialingWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.subscription.trialing.webhook") //struct name
	event.SetType(EventCreemSubscriptionTrialingWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
