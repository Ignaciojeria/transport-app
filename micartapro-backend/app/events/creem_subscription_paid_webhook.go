package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemSubscriptionPaidWebhook struct {
	ID        string                       `json:"id" example:"evt_21mO1jWmU2QHe7u2oFV7y1"`
	EventType string                       `json:"eventType" example:"subscription.paid"`
	CreatedAt int64                        `json:"created_at" example:"1728734327355"`
	Object    CreemSubscriptionPaidObject  `json:"object"`
}

type CreemSubscriptionPaidProduct struct {
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

type CreemSubscriptionPaidCustomer struct {
	ID        string `json:"id" example:"cust_1OcIK1GEuVvXZwD19tjq2z"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Tester Test"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2024-10-11T09:16:48.557Z"`
	UpdatedAt string `json:"updated_at" example:"2024-10-11T09:16:48.557Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemSubscriptionPaidObject struct {
	ID                      string                          `json:"id" example:"sub_6pC2lNB6joCRQIZ1aMrTpi"`
	Object                  string                          `json:"object" example:"subscription"`
	Product                 CreemSubscriptionPaidProduct    `json:"product"`
	Customer                CreemSubscriptionPaidCustomer   `json:"customer"`
	CollectionMethod        string                          `json:"collection_method" example:"charge_automatically"`
	Status                  string                          `json:"status" example:"active"`
	LastTransactionID       string                          `json:"last_transaction_id" example:"tran_5yMaWzAl3jxuGJMCOrYWwk"`
	LastTransactionDate     string                          `json:"last_transaction_date" example:"2024-10-12T11:58:47.109Z"`
	NextTransactionDate     string                          `json:"next_transaction_date" example:"2024-11-12T11:58:38.000Z"`
	CurrentPeriodStartDate  string                          `json:"current_period_start_date" example:"2024-10-12T11:58:38.000Z"`
	CurrentPeriodEndDate    string                          `json:"current_period_end_date" example:"2024-11-12T11:58:38.000Z"`
	CanceledAt              any                             `json:"canceled_at" example:"null"`
	CreatedAt               string                          `json:"created_at" example:"2024-10-12T11:58:45.425Z"`
	UpdatedAt               string                          `json:"updated_at" example:"2024-10-12T11:58:45.425Z"`
	Metadata                CreemSubscriptionTrialingMetadata `json:"metadata"`
	Mode                    string                          `json:"mode" example:"local"`
}

func (c CreemSubscriptionPaidWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.subscription.paid.webhook")
	event.SetType(EventCreemSubscriptionPaidWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
