package mercadopago

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-resty/resty/v2"
)

type GetMercadoPagoPayment func(ctx context.Context, paymentID string) (PaymentResponse, error)

type PaymentResponse struct {
	ID                 int64                  `json:"id"`
	Status             string                 `json:"status"`
	StatusDetail       string                 `json:"status_detail"`
	ExternalReference  string                 `json:"external_reference,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
	DateCreated        string                 `json:"date_created"`
	DateLastUpdated    string                 `json:"date_last_updated"`
	TransactionAmount  float64                `json:"transaction_amount"`
	CurrencyID         string                 `json:"currency_id"`
	PaymentMethodID    string                 `json:"payment_method_id"`
	PaymentTypeID      string                 `json:"payment_type_id"`
	TransactionDetails *TransactionDetails    `json:"transaction_details,omitempty"`
}

type TransactionDetails struct {
	TotalPaidAmount   float64 `json:"total_paid_amount"`
	NetReceivedAmount float64 `json:"net_received_amount"`
}

func init() {
	ioc.Register(NewGetMercadoPagoPayment)
}

func NewGetMercadoPagoPayment(cli *resty.Client, conf configuration.Conf) GetMercadoPagoPayment {
	return func(ctx context.Context, paymentID string) (PaymentResponse, error) {
		url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentID)

		resp, err := cli.R().
			SetContext(ctx).
			SetHeader("Authorization", "Bearer "+conf.MERCADOPAGO_ACCESS_TOKEN).
			Get(url)

		if err != nil {
			return PaymentResponse{}, fmt.Errorf("failed to make HTTP request: %w", err)
		}

		if resp.IsError() {
			return PaymentResponse{}, fmt.Errorf("Mercado Pago API error (status %d): %s", resp.StatusCode(), resp.String())
		}

		var result PaymentResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return PaymentResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return result, nil
	}
}
