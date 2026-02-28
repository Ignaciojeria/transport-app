package mercadopago

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-resty/resty/v2"
)

type GetMercadoPagoOrder func(ctx context.Context, orderID string) (OrderResponse, error)

type OrderResponse struct {
	ID                string                 `json:"id"`
	Status            string                 `json:"status"`
	StatusDetail      string                 `json:"status_detail"`
	ExternalReference string                 `json:"external_reference,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	DateCreated       string                 `json:"date_created"`
	DateLastUpdated   string                 `json:"date_last_updated"`
	Items             []OrderItem            `json:"items,omitempty"`
	TotalAmount       float64                `json:"total_amount,omitempty"`
}

func init() {
	ioc.Register(NewGetMercadoPagoOrder)
}

func NewGetMercadoPagoOrder(cli *resty.Client, conf configuration.Conf) GetMercadoPagoOrder {
	return func(ctx context.Context, orderID string) (OrderResponse, error) {
		url := fmt.Sprintf("https://api.mercadopago.com/checkout/preferences/%s", orderID)

		resp, err := cli.R().
			SetContext(ctx).
			SetHeader("Authorization", "Bearer "+conf.MERCADOPAGO_ACCESS_TOKEN).
			Get(url)

		if err != nil {
			return OrderResponse{}, fmt.Errorf("failed to make HTTP request: %w", err)
		}

		if resp.IsError() {
			return OrderResponse{}, fmt.Errorf("Mercado Pago API error (status %d): %s", resp.StatusCode(), resp.String())
		}

		var result OrderResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return OrderResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return result, nil
	}
}
