package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemRefundCreatedWebhook struct {
	ID        string                    `json:"id" example:"evt_61eTsJHUgInFw2BQKhTiPV"`
	EventType string                    `json:"eventType" example:"refund.created"`
	CreatedAt int64                     `json:"created_at" example:"1728734351631"`
	Object    CreemRefundCreatedObject  `json:"object"`
}

type CreemRefundCreatedTransaction struct {
	ID            string `json:"id" example:"tran_5yMaWzAl3jxuGJMCOrYWwk"`
	Object        string `json:"object" example:"transaction"`
	Amount        int    `json:"amount" example:"1000"`
	AmountPaid    int    `json:"amount_paid" example:"1210"`
	Currency      string `json:"currency" example:"EUR"`
	Type          string `json:"type" example:"invoice"`
	TaxCountry    string `json:"tax_country" example:"NL"`
	TaxAmount     int    `json:"tax_amount" example:"210"`
	Status        string `json:"status" example:"refunded"`
	RefundedAmount int   `json:"refunded_amount" example:"1210"`
	Order         string `json:"order" example:"ord_4aDwWXjMLpes4Kj4XqNnUA"`
	Subscription  string `json:"subscription" example:"sub_6pC2lNB6joCRQIZ1aMrTpi"`
	Description   string `json:"description" example:"Subscription payment"`
	PeriodStart   int64  `json:"period_start" example:"1728734318000"`
	PeriodEnd     int64  `json:"period_end" example:"1731412718000"`
	CreatedAt     int64  `json:"created_at" example:"1728734327109"`
	Mode          string `json:"mode" example:"local"`
}

type CreemRefundCreatedSubscription struct {
	ID                     string                            `json:"id" example:"sub_6pC2lNB6joCRQIZ1aMrTpi"`
	Object                 string                            `json:"object" example:"subscription"`
	Product                string                            `json:"product" example:"prod_d1AY2Sadk9YAvLI0pj97f"`
	Customer               string                            `json:"customer" example:"cust_1OcIK1GEuVvXZwD19tjq2z"`
	CollectionMethod       string                            `json:"collection_method" example:"charge_automatically"`
	Status                 string                            `json:"status" example:"canceled"`
	LastTransactionID      string                            `json:"last_transaction_id" example:"tran_5yMaWzAl3jxuGJMCOrYWwk"`
	LastTransactionDate    string                            `json:"last_transaction_date" example:"2024-10-12T11:58:47.109Z"`
	CurrentPeriodStartDate string                            `json:"current_period_start_date" example:"2024-10-12T11:58:38.000Z"`
	CurrentPeriodEndDate   string                            `json:"current_period_end_date" example:"2024-11-12T11:58:38.000Z"`
	CanceledAt             string                            `json:"canceled_at" example:"2024-10-12T11:58:57.813Z"`
	CreatedAt              string                            `json:"created_at" example:"2024-10-12T11:58:45.425Z"`
	UpdatedAt              string                            `json:"updated_at" example:"2024-10-12T11:58:57.827Z"`
	Metadata               CreemSubscriptionTrialingMetadata  `json:"metadata"`
	Mode                   string                            `json:"mode" example:"local"`
}

type CreemRefundCreatedCheckout struct {
	ID           string                            `json:"id" example:"ch_4l0N34kxo16AhRKUHFUuXr"`
	Object       string                            `json:"object" example:"checkout"`
	RequestID    string                            `json:"request_id" example:"my-request-id"`
	CustomFields []interface{}                     `json:"custom_fields"`
	Status       string                            `json:"status" example:"completed"`
	Metadata     CreemSubscriptionTrialingMetadata `json:"metadata"`
	Mode         string                            `json:"mode" example:"local"`
}

type CreemRefundCreatedOrder struct {
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

type CreemRefundCreatedCustomer struct {
	ID        string `json:"id" example:"cust_1OcIK1GEuVvXZwD19tjq2z"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Tester Test"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2024-10-11T09:16:48.557Z"`
	UpdatedAt string `json:"updated_at" example:"2024-10-11T09:16:48.557Z"`
	Mode      string `json:"mode" example:"local"`
}

type CreemRefundCreatedObject struct {
	ID            string                        `json:"id" example:"ref_3DB9NQFvk18TJwSqd0N6bd"`
	Object        string                        `json:"object" example:"refund"`
	Status        string                        `json:"status" example:"succeeded"`
	RefundAmount  int                           `json:"refund_amount" example:"1210"`
	RefundCurrency string                       `json:"refund_currency" example:"EUR"`
	Reason        string                        `json:"reason" example:"requested_by_customer"`
	Transaction   CreemRefundCreatedTransaction `json:"transaction"`
	Subscription  CreemRefundCreatedSubscription `json:"subscription"`
	Checkout      CreemRefundCreatedCheckout     `json:"checkout"`
	Order         CreemRefundCreatedOrder       `json:"order"`
	Customer      CreemRefundCreatedCustomer     `json:"customer"`
	CreatedAt     int64                         `json:"created_at" example:"1728734351525"`
	Mode          string                        `json:"mode" example:"local"`
}

func (c CreemRefundCreatedWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.refund.created.webhook")
	event.SetType(EventCreemRefundCreatedWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
