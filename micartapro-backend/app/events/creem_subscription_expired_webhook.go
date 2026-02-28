package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemSubscriptionExpiredWebhook struct {
	ID        string                         `json:"id" example:"evt_V5CxhipUu10BYonO2Vshb"`
	EventType string                         `json:"eventType" example:"subscription.expired"`
	CreatedAt int64                          `json:"created_at" example:"1734463872058"`
	Object    CreemSubscriptionExpiredObject `json:"object"`
}

type CreemSubscriptionExpiredProduct struct {
	ID                string `json:"id" example:"prod_3ELsC3Lt97orn81SOdgQI3"`
	Name              string `json:"name" example:"Subs"`
	Description       string `json:"description" example:"Subs"`
	ImageURL          any    `json:"image_url" example:"null"`
	Price             int    `json:"price" example:"1200"`
	Currency          string `json:"currency" example:"EUR"`
	BillingType       string `json:"billing_type" example:"recurring"`
	BillingPeriod     string `json:"billing_period" example:"every-year"`
	Status            string `json:"status" example:"active"`
	TaxMode           string `json:"tax_mode" example:"exclusive"`
	TaxCategory       string `json:"tax_category" example:"saas"`
	DefaultSuccessURL string `json:"default_success_url" example:""`
	CreatedAt         string `json:"created_at" example:"2024-12-11T17:33:32.186Z"`
	UpdatedAt         string `json:"updated_at" example:"2024-12-11T17:33:32.186Z"`
	Mode              string `json:"mode" example:"local"`
}

type CreemSubscriptionExpiredCustomer struct {
	ID        string `json:"id" example:"cust_3y4k2CELGsw7n9Eeeiw2hm"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Alec Erasmus"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2024-12-09T16:09:20.709Z"`
	UpdatedAt string `json:"updated_at" example:"2024-12-09T16:09:20.709Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemSubscriptionExpiredObject struct {
	ID                     string                           `json:"id" example:"sub_7FgHvrOMC28tG5DEemoCli"`
	Object                 string                           `json:"object" example:"subscription"`
	Product                CreemSubscriptionExpiredProduct  `json:"product"`
	Customer               CreemSubscriptionExpiredCustomer `json:"customer"`
	CollectionMethod       string                           `json:"collection_method" example:"charge_automatically"`
	Status                 string                           `json:"status" example:"active"`
	LastTransactionID      string                           `json:"last_transaction_id" example:"tran_6ZeTvMqMkGdAIIjw5aAcnh"`
	LastTransactionDate    string                           `json:"last_transaction_date" example:"2024-12-16T12:40:12.658Z"`
	NextTransactionDate    string                           `json:"next_transaction_date" example:"2025-12-16T12:39:47.000Z"`
	CurrentPeriodStartDate string                           `json:"current_period_start_date" example:"2024-12-16T12:39:47.000Z"`
	CurrentPeriodEndDate   string                           `json:"current_period_end_date" example:"2024-12-16T12:39:47.000Z"`
	CanceledAt             any                              `json:"canceled_at" example:"null"`
	CreatedAt              string                           `json:"created_at" example:"2024-12-16T12:40:05.058Z"`
	UpdatedAt              string                           `json:"updated_at" example:"2024-12-16T12:40:05.058Z"`
	Mode                   string                           `json:"mode" example:"local"`
}

func (c CreemSubscriptionExpiredWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.subscription.expired.webhook")
	event.SetType(EventCreemSubscriptionExpiredWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
