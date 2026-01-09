package events

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CreemDisputeCreatedWebhook struct {
	ID        string                      `json:"id" example:"evt_6mfLDL7P0NYwYQqCrICvDH"`
	EventType string                      `json:"eventType" example:"dispute.created"`
	CreatedAt int64                       `json:"created_at" example:"1750941264812"`
	Object    CreemDisputeCreatedObject   `json:"object"`
}

type CreemDisputeCreatedTransaction struct {
	ID            string `json:"id" example:"tran_4Dk8CxWFdceRUQgMFhCCXX"`
	Object        string `json:"object" example:"transaction"`
	Amount        int    `json:"amount" example:"1100"`
	AmountPaid    int    `json:"amount_paid" example:"1331"`
	Currency      string `json:"currency" example:"EUR"`
	Type          string `json:"type" example:"invoice"`
	TaxCountry    string `json:"tax_country" example:"NL"`
	TaxAmount     int    `json:"tax_amount" example:"231"`
	Status        string `json:"status" example:"chargeback"`
	RefundedAmount int   `json:"refunded_amount" example:"1331"`
	Order         string `json:"order" example:"ord_57bf8042UmG8fFypxZrfnj"`
	Subscription  string `json:"subscription" example:"sub_5sD6zM482uwOaEoyEUDDJs"`
	Customer      string `json:"customer" example:"cust_OJPZd2GMxgo1MGPNXXBSN"`
	Description   string `json:"description" example:"Subscription payment"`
	PeriodStart   int64  `json:"period_start" example:"1750941201000"`
	PeriodEnd     int64  `json:"period_end" example:"1753533201000"`
	CreatedAt     int64  `json:"created_at" example:"1750941205659"`
	Mode          string `json:"mode" example:"sandbox"`
}

type CreemDisputeCreatedSubscription struct {
	ID                     string `json:"id" example:"sub_5sD6zM482uwOaEoyEUDDJs"`
	Object                 string `json:"object" example:"subscription"`
	Product                string `json:"product" example:"prod_3EFtQRQ9SNIizK3xwfxZHu"`
	Customer               string `json:"customer" example:"cust_OJPZd2GMxgo1MGPNXXBSN"`
	CollectionMethod       string `json:"collection_method" example:"charge_automatically"`
	Status                 string `json:"status" example:"active"`
	CurrentPeriodStartDate string `json:"current_period_start_date" example:"2025-06-26T12:33:21.000Z"`
	CurrentPeriodEndDate   string `json:"current_period_end_date" example:"2025-07-26T12:33:21.000Z"`
	CanceledAt             any    `json:"canceled_at" example:"null"`
	CreatedAt              string `json:"created_at" example:"2025-06-26T12:33:23.589Z"`
	UpdatedAt              string `json:"updated_at" example:"2025-06-26T12:33:26.102Z"`
	Mode                   string `json:"mode" example:"sandbox"`
}

type CreemDisputeCreatedCheckoutCustomField struct {
	Key       string `json:"key" example:"testing"`
	Text      map[string]interface{} `json:"text"`
	Type      string `json:"type" example:"text"`
	Label     string `json:"label" example:"Testing"`
	Optional  bool   `json:"optional" example:"false"`
}

type CreemDisputeCreatedCheckout struct {
	ID           string                              `json:"id" example:"ch_1bJMvqGGzHIftf4ewLXJeq"`
	Object       string                              `json:"object" example:"checkout"`
	Product      string                              `json:"product" example:"prod_3EFtQRQ9SNIizK3xwfxZHu"`
	Units        int                                 `json:"units" example:"1"`
	CustomFields []CreemDisputeCreatedCheckoutCustomField `json:"custom_fields"`
	Status       string                              `json:"status" example:"completed"`
	Mode         string                              `json:"mode" example:"sandbox"`
}

type CreemDisputeCreatedOrder struct {
	Object      string `json:"object" example:"order"`
	ID          string `json:"id" example:"ord_57bf8042UmG8fFypxZrfnj"`
	Customer    string `json:"customer" example:"cust_OJPZd2GMxgo1MGPNXXBSN"`
	Product     string `json:"product" example:"prod_3EFtQRQ9SNIizK3xwfxZHu"`
	Amount      int    `json:"amount" example:"1100"`
	Currency    string `json:"currency" example:"EUR"`
	SubTotal    int    `json:"sub_total" example:"1100"`
	TaxAmount   int    `json:"tax_amount" example:"231"`
	AmountDue   int    `json:"amount_due" example:"1331"`
	AmountPaid  int    `json:"amount_paid" example:"1331"`
	Status      string `json:"status" example:"paid"`
	Type        string `json:"type" example:"recurring"`
	Transaction string `json:"transaction" example:"tran_4Dk8CxWFdceRUQgMFhCCXX"`
	CreatedAt   string `json:"created_at" example:"2025-06-26T12:32:41.395Z"`
	UpdatedAt   string `json:"updated_at" example:"2025-06-26T12:32:41.395Z"`
	Mode        string `json:"mode" example:"sandbox"`
}

type CreemDisputeCreatedCustomer struct {
	ID        string `json:"id" example:"cust_OJPZd2GMxgo1MGPNXXBSN"`
	Object    string `json:"object" example:"customer"`
	Email     string `json:"email" example:"customer@emaildomain"`
	Name      string `json:"name" example:"Alec Erasmus"`
	Country   string `json:"country" example:"NL"`
	CreatedAt string `json:"created_at" example:"2025-02-05T10:11:01.146Z"`
	UpdatedAt string `json:"updated_at" example:"2025-02-05T10:11:01.146Z"`
	Mode      string `json:"mode" example:"sandbox"`
}

type CreemDisputeCreatedObject struct {
	ID           string                        `json:"id" example:"disp_6vSsOdTANP5PhOzuDlUuXE"`
	Object       string                        `json:"object" example:"dispute"`
	Amount       int                           `json:"amount" example:"1331"`
	Currency     string                        `json:"currency" example:"EUR"`
	Transaction  CreemDisputeCreatedTransaction `json:"transaction"`
	Subscription CreemDisputeCreatedSubscription `json:"subscription"`
	Checkout     CreemDisputeCreatedCheckout    `json:"checkout"`
	Order        CreemDisputeCreatedOrder       `json:"order"`
	Customer     CreemDisputeCreatedCustomer    `json:"customer"`
	CreatedAt    int64                         `json:"created_at" example:"1750941264728"`
	Mode         string                         `json:"mode" example:"local"`
}

func (c CreemDisputeCreatedWebhook) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("creem.dispute.created.webhook")
	event.SetType(EventCreemDisputeCreatedWebhook)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}
