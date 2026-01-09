package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemCheckoutCompletedWebhook struct {
	ID        string                       `json:"id" example:"evt_5WHHcZPv7VS0YUsberIuOz"`
	EventType string                       `json:"eventType" example:"checkout.completed"`
	CreatedAt int64                        `json:"created_at" example:"1728734325927"`
	Object    CreemCheckoutCompletedObject `json:"object"`
}

type CreemCheckoutCompletedProduct struct {
	ID                string `json:"id" example:"prod_d1AY2Sadk9YAvLI0pj97f"`
	Name              string `json:"name" example:"Monthly"`
	Description       string `json:"description" example:"Monthly"`
	ImageURL          any    `json:"image_url" example:"null"`
	Price             int    `json:"price" example:"1000"`
	Currency          string `json:"currency" example:"EUR"`
	BillingType       string `json:"billing_type" example:"recurring"`
	BillingPeriod     string `json:"billing_period" example:"every-month"`
	Status            string `json:"status" example:"active"`
	TaxMode           string `json:"tax_mode" example:"exclusive"`
	TaxCategory       string `json:"tax_category" example:"saas"`
	DefaultSuccessURL string `json:"default_success_url" example:""`
	CreatedAt         string `json:"created_at" example:"2024-10-11T11:50:00.182Z"`
	UpdatedAt         string `json:"updated_at" example:"2024-10-11T11:50:00.182Z"`
	Mode              string `json:"mode" example:"local"`
}

type CreemCheckoutCompletedCustomer struct {
	ID        string `json:"id" example:"cust_1OcIK1GEuVvXZwD19tjq2z"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Tester Test"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2024-10-11T09:16:48.557Z"`
	UpdatedAt string `json:"updated_at" example:"2024-10-11T09:16:48.557Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemCheckoutCompletedOrder struct {
	ID        string `json:"id" example:"ord_4aDwWXjMLpes4Kj4XqNnUA"`
	Customer  string `json:"customer" example:"cust_1OcIK1GEuVvXZwD19tjq2z"`
	Product   string `json:"product" example:"prod_d1AY2Sadk9YAvLI0pj97f"`
	Amount    int    `json:"amount" example:"1000"`
	Currency  string `json:"currency" example:"EUR"`
	Status    string `json:"status" example:"paid"`
	Type      string `json:"type" example:"recurring"`
	CreatedAt string `json:"created_at" example:"2024-10-12T11:58:33.097Z"`
	UpdatedAt string `json:"updated_at" example:"2024-10-12T11:58:33.097Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemCheckoutCompletedSubscription struct {
	ID               string                            `json:"id" example:"sub_6pC2lNB6joCRQIZ1aMrTpi"`
	Object           string                            `json:"object" example:"subscription"`
	Product          string                            `json:"product" example:"prod_d1AY2Sadk9YAvLI0pj97f"`
	Customer         string                            `json:"customer" example:"cust_1OcIK1GEuVvXZwD19tjq2z"`
	CollectionMethod string                            `json:"collection_method" example:"charge_automatically"`
	Status           string                            `json:"status" example:"active"`
	CanceledAt       any                               `json:"canceled_at" example:"null"`
	CreatedAt        string                            `json:"created_at" example:"2024-10-12T11:58:45.425Z"`
	UpdatedAt        string                            `json:"updated_at" example:"2024-10-12T11:58:45.425Z"`
	Metadata         CreemSubscriptionTrialingMetadata `json:"metadata"`
	Mode             string                            `json:"mode" example:"local"`
}

type CreemCheckoutCompletedObject struct {
	ID           string                             `json:"id" example:"ch_4l0N34kxo16AhRKUHFUuXr"`
	Object       string                             `json:"object" example:"checkout"`
	RequestID    string                             `json:"request_id" example:"my-request-id"`
	Order        CreemCheckoutCompletedOrder        `json:"order"`
	Product      CreemCheckoutCompletedProduct      `json:"product"`
	Customer     CreemCheckoutCompletedCustomer     `json:"customer"`
	Subscription CreemCheckoutCompletedSubscription `json:"subscription"`
	CustomFields []interface{}                      `json:"custom_fields"`
	Status       string                             `json:"status" example:"completed"`
	Metadata     CreemSubscriptionTrialingMetadata  `json:"metadata"`
	Mode         string                             `json:"mode" example:"local"`
}

func (c CreemCheckoutCompletedWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.checkout.completed.webhook")
	event.SetType(EventCreemCheckoutCompletedWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
