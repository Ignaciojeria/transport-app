package restyclient

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/httpresty"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type GetCreemCheckoutUrl func(ctx context.Context, userID string) (creemCheckoutResponse, error)

type creemCheckoutRequest struct {
	ProductID  string            `json:"product_id"`
	SuccessURL string            `json:"success_url"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type creemCheckoutResponse struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Product     string `json:"product"`
	Units       int    `json:"units"`
	Status      string `json:"status"`
	CheckoutURL string `json:"checkout_url"`
	SuccessURL  string `json:"success_url"`
	Mode        string `json:"mode"`
}

func init() {
	ioc.Registry(NewGetCreemCheckoutUrl, httpresty.NewClient, configuration.NewConf)
}
func NewGetCreemCheckoutUrl(cli *resty.Client, conf configuration.Conf) GetCreemCheckoutUrl {
	return func(ctx context.Context, userID string) (creemCheckoutResponse, error) {
		url := conf.CREEM_DNS + "/v1/checkouts"

		req := creemCheckoutRequest{
			ProductID:  conf.CREEM_PRODUCT_ID,
			SuccessURL: conf.CREEM_SUCCESS_URL + "?payment=success",
			Metadata: map[string]string{
				"user_id": userID,
			},
		}

		resp, err := cli.R().
			SetContext(ctx).
			SetHeader("x-api-key", conf.CREEM_API_KEY).
			SetBody(req).
			Post(url)

		if err != nil {
			return creemCheckoutResponse{}, fmt.Errorf("failed to make HTTP request: %w", err)
		}

		if resp.IsError() {
			return creemCheckoutResponse{}, fmt.Errorf("API error (status %d): %s", resp.StatusCode(), resp.String())
		}

		var result creemCheckoutResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return creemCheckoutResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if result.CheckoutURL == "" {
			return creemCheckoutResponse{}, fmt.Errorf("checkout_url not found in response")
		}

		return result, nil
	}
}
