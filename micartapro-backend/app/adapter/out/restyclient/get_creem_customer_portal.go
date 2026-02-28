package restyclient

import (
	"context"
	"encoding/json"
	"fmt"
	"micartapro/app/shared/configuration"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-resty/resty/v2"
)

type GetCreemCustomerPortal func(ctx context.Context, customerID string) (creemCustomerPortalResponse, error)

type creemCustomerPortalRequest struct {
	CustomerID string `json:"customer_id"`
}

type creemCustomerPortalResponse struct {
	CustomerPortalLink string `json:"customer_portal_link"`
}

func init() {
	ioc.Register(NewGetCreemCustomerPortal)
}

func NewGetCreemCustomerPortal(cli *resty.Client, conf configuration.Conf) GetCreemCustomerPortal {
	return func(ctx context.Context, customerID string) (creemCustomerPortalResponse, error) {
		url := conf.CREEM_DNS + "/v1/customers/billing"

		req := creemCustomerPortalRequest{
			CustomerID: customerID,
		}

		resp, err := cli.R().
			SetContext(ctx).
			SetHeader("x-api-key", conf.CREEM_API_KEY).
			SetHeader("Content-Type", "application/json").
			SetBody(req).
			Post(url)

		if err != nil {
			return creemCustomerPortalResponse{}, fmt.Errorf("failed to make HTTP request: %w", err)
		}

		if resp.IsError() {
			return creemCustomerPortalResponse{}, fmt.Errorf("API error (status %d): %s", resp.StatusCode(), resp.String())
		}

		var result creemCustomerPortalResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return creemCustomerPortalResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if result.CustomerPortalLink == "" {
			return creemCustomerPortalResponse{}, fmt.Errorf("customer_portal_link not found in response")
		}

		return result, nil
	}
}
