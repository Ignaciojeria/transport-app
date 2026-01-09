package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemSubscriptionUpdateWebhook struct {
	ID        string                        `json:"id" example:"evt_5pJMUuvqaqvttFVUvtpY32"`
	EventType string                        `json:"eventType" example:"subscription.update"`
	CreatedAt int64                         `json:"created_at" example:"1737890536421"`
	Object    CreemSubscriptionUpdateObject `json:"object"`
}

type CreemSubscriptionUpdateProduct struct {
	ID                string `json:"id" example:"prod_1dP15yoyogQe2seEt1Evf3"`
	Name              string `json:"name" example:"Monthly Sub"`
	Description       string `json:"description" example:"Test Test"`
	ImageURL          any    `json:"image_url" example:"null"`
	Price             int    `json:"price" example:"1000"`
	Currency          string `json:"currency" example:"EUR"`
	BillingType       string `json:"billing_type" example:"recurring"`
	BillingPeriod     string `json:"billing_period" example:"every-month"`
	Status            string `json:"status" example:"active"`
	TaxMode           string `json:"tax_mode" example:"exclusive"`
	TaxCategory       string `json:"tax_category" example:"saas"`
	DefaultSuccessURL string `json:"default_success_url" example:""`
	CreatedAt         string `json:"created_at" example:"2025-01-26T11:17:16.082Z"`
	UpdatedAt         string `json:"updated_at" example:"2025-01-26T11:17:16.082Z"`
	Mode              string `json:"mode" example:"local"`
}

type CreemSubscriptionUpdateCustomer struct {
	ID        string `json:"id" example:"cust_2fQZKKUZqtNhH2oDWevQkW"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"John Doe"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2025-01-26T11:18:24.071Z"`
	UpdatedAt string `json:"updated_at" example:"2025-01-26T11:18:24.071Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemSubscriptionUpdateItems struct {
	Object    string `json:"object" example:"subscription_item"`
	ID        string `json:"id" example:"sitem_3QWlqRbAat2eBRakAxFtt9"`
	ProductID string `json:"product_id" example:"prod_5jnudVkLGZWF4AqMFBs5t5"`
	PriceID   string `json:"price_id" example:"pprice_4W0mJK6uGiQzHbVhfaFTl1"`
	Units     int    `json:"units" example:"1"`
	CreatedAt string `json:"created_at" example:"2025-01-26T11:20:40.296Z"`
	UpdatedAt string `json:"updated_at" example:"2025-01-26T11:20:40.296Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemSubscriptionUpdateObject struct {
	ID                     string                          `json:"id" example:"sub_2qAuJgWmXhXHAuef9k4Kur"`
	Object                 string                          `json:"object" example:"subscription"`
	Product                CreemSubscriptionUpdateProduct  `json:"product"`
	Customer               CreemSubscriptionUpdateCustomer `json:"customer"`
	Items                  []CreemSubscriptionUpdateItems  `json:"items"`
	CollectionMethod       string                          `json:"collection_method" example:"charge_automatically"`
	Status                 string                          `json:"status" example:"active"`
	CurrentPeriodStartDate string                          `json:"current_period_start_date" example:"2025-01-26T11:20:36.000Z"`
	CurrentPeriodEndDate   string                          `json:"current_period_end_date" example:"2025-02-26T11:20:36.000Z"`
	CanceledAt             any                             `json:"canceled_at" example:"null"`
	CreatedAt              string                          `json:"created_at" example:"2025-01-26T11:20:40.292Z"`
	UpdatedAt              string                          `json:"updated_at" example:"2025-01-26T11:22:16.388Z"`
	Mode                   string                          `json:"mode" example:"local"`
}

func (c CreemSubscriptionUpdateWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.subscription.update.webhook")
	event.SetType(EventCreemSubscriptionUpdateWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
