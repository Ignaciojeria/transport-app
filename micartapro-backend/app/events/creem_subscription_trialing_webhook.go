package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemSubscriptionTrialingWebhook struct {
	ID        string                          `json:"id" example:"evt_2ciAM8ABYtj0pVueeJPxUZ"`
	EventType string                          `json:"eventType" example:"subscription.trialing"`
	CreatedAt int64                           `json:"created_at" example:"1739963911073"`
	Object    CreemSubscriptionTrialingObject `json:"object"`
}

type CreemSubscriptionTrialingProduct struct {
	ID                string `json:"id" example:"prod_3kpf0ZdpcfsSCQ3kDiwg9m"`
	Name              string `json:"name" example:"trail"`
	Description       string `json:"description" example:"asdfasf"`
	ImageURL          any    `json:"image_url" example:"null"`
	Price             int    `json:"price" example:"1100"`
	Currency          string `json:"currency" example:"EUR"`
	BillingType       string `json:"billing_type" example:"recurring"`
	BillingPeriod     string `json:"billing_period" example:"every-month"`
	Status            string `json:"status" example:"active"`
	TaxMode           string `json:"tax_mode" example:"exclusive"`
	TaxCategory       string `json:"tax_category" example:"saas"`
	DefaultSuccessURL string `json:"default_success_url" example:""`
	CreatedAt         string `json:"created_at" example:"2025-02-19T11:18:07.570Z"`
	UpdatedAt         string `json:"updated_at" example:"2025-02-19T11:18:07.570Z"`
	Mode              string `json:"mode" example:"test"`
}

type CreemSubscriptionTrialingCustomer struct {
	ID        string `json:"id" example:"cust_4fpU8kYkQmI1XKBwU2qeME"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Alec Erasmus"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2024-11-07T23:21:11.763Z"`
	UpdatedAt string `json:"updated_at" example:"2024-11-07T23:21:11.763Z"`
	Mode      string `json:"mode" example:"test"`
}

type CreemSubscriptionTrialingItems struct {
	Object    string `json:"object" example:"subscription_item"`
	ID        string `json:"id" example:"sitem_1xbHCmIM61DHGRBCFn0W1L"`
	ProductID string `json:"product_id" example:"prod_3kpf0ZdpcfsSCQ3kDiwg9m"`
	PriceID   string `json:"price_id" example:"pprice_517h9CebmM3P079bGAXHnE"`
	Units     int    `json:"units" example:"1"`
	CreatedAt string `json:"created_at" example:"2025-02-19T11:18:30.690Z"`
	UpdatedAt string `json:"updated_at" example:"2025-02-19T11:18:30.690Z"`
	Mode      string `json:"mode" example:"test"`
}

type CreemSubscriptionTrialingMetadata struct {
	UserID string `json:"user_id" example:"763a590a-9b8e-4a91-b8ee-47f2a64d003d"`
}

type CreemSubscriptionTrialingObject struct {
	ID                     string                            `json:"id" example:"sub_dxiauR8zZOwULx5QM70wJ"`
	Object                 string                            `json:"object" example:"subscription"`
	Product                CreemSubscriptionTrialingProduct  `json:"product"`
	Customer               CreemSubscriptionTrialingCustomer `json:"customer"`
	Items                  []CreemSubscriptionTrialingItems  `json:"items"`
	CollectionMethod       string                            `json:"collection_method" example:"charge_automatically"`
	Status                 string                            `json:"status" example:"trialing"`
	CurrentPeriodStartDate string                            `json:"current_period_start_date" example:"2025-02-19T11:18:25.000Z"`
	CurrentPeriodEndDate   string                            `json:"current_period_end_date" example:"2025-02-26T11:18:25.000Z"`
	CanceledAt             any                               `json:"canceled_at" example:"null"`
	CreatedAt              string                            `json:"created_at" example:"2025-02-19T11:18:30.674Z"`
	UpdatedAt              string                            `json:"updated_at" example:"2025-02-19T11:18:30.674Z"`
	Metadata               CreemSubscriptionTrialingMetadata `json:"metadata"`
	Mode                   string                            `json:"mode" example:"test"`
}

func (c CreemSubscriptionTrialingWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.subscription.trialing.webhook")
	event.SetType(EventCreemSubscriptionTrialingWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
