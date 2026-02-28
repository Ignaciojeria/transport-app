package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemSubscriptionActiveWebhook struct {
	ID        string                        `json:"id" example:"evt_6EptlmjazyGhEPiNQ5f4lz"`
	EventType string                        `json:"eventType" example:"subscription.active"`
	CreatedAt int64                         `json:"created_at" example:"1728734325927"`
	Object    CreemSubscriptionActiveObject `json:"object"`
}

type CreemSubscriptionActiveProduct struct {
	ID                string `json:"id" example:"prod_AnVJ11ujp7x953ARpJvAF"`
	Name              string `json:"name" example:"My Product - Product 01"`
	Description       string `json:"description" example:"Test my product"`
	ImageURL          any    `json:"image_url" example:"null"`
	Price             int    `json:"price" example:"10000"`
	Currency          string `json:"currency" example:"EUR"`
	BillingType       string `json:"billing_type" example:"recurring"`
	BillingPeriod     string `json:"billing_period" example:"every-month"`
	Status            string `json:"status" example:"active"`
	TaxMode           string `json:"tax_mode" example:"inclusive"`
	TaxCategory       string `json:"tax_category" example:"saas"`
	DefaultSuccessURL string `json:"default_success_url" example:""`
	CreatedAt         string `json:"created_at" example:"2024-09-16T16:12:09.813Z"`
	UpdatedAt         string `json:"updated_at" example:"2024-09-16T16:12:09.813Z"`
	Mode              string `json:"mode" example:"local"`
}

type CreemSubscriptionActiveCustomer struct {
	ID        string `json:"id" example:"cust_3biFPNt4Cz5YRDSdIqs7kc"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Tester Test"`
	Country   string `json:"country" example:"SE"`
	CreatedAt string `json:"created_at" example:"2024-09-16T16:13:39.265Z"`
	UpdatedAt string `json:"updated_at" example:"2024-09-16T16:13:39.265Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemSubscriptionActiveObject struct {
	ID               string                          `json:"id" example:"sub_21lfZb67szyvMiXnm6SVi0"`
	Object           string                          `json:"object" example:"subscription"`
	Product          CreemSubscriptionActiveProduct  `json:"product"`
	Customer         CreemSubscriptionActiveCustomer `json:"customer"`
	CollectionMethod string                          `json:"collection_method" example:"charge_automatically"`
	Status           string                          `json:"status" example:"active"`
	CanceledAt       string                          `json:"canceled_at" example:"2024-09-16T19:40:41.984Z"`
	CreatedAt        string                          `json:"created_at" example:"2024-09-16T19:40:41.984Z"`
	UpdatedAt        string                          `json:"updated_at" example:"2024-09-16T19:40:42.121Z"`
	Mode             string                          `json:"mode" example:"local"`
}

func (c CreemSubscriptionActiveWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.subscription.active.webhook")
	event.SetType(EventCreemSubscriptionActiveWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
