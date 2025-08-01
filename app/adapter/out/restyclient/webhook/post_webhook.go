package webhook

import (
	"context"
	"fmt"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/httpresty"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type PostWebhook func(ctx context.Context, input interface{}, webhookType string) error

func init() {
	ioc.Registry(NewPostWebhook, httpresty.NewClient, configuration.NewConf)
}

func NewPostWebhook(c *resty.Client, config configuration.Conf) PostWebhook {
	return func(ctx context.Context, input interface{}, webhookType string) error {
		accessToken, ok := sharedcontext.AccessTokenFromContext(ctx)
		if !ok {
			return fmt.Errorf("access token not found in context")
		}
		// Crear la petición HTTP con resty
		resp, err := c.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeaders(map[string]string{
				"X-Access-Token": accessToken,
				"tenant":         sharedcontext.TenantIDFromContext(ctx).String() + "-" + sharedcontext.TenantCountryFromContext(ctx),
			}).
			SetBody(input).
			Post(config.MASTER_NODE_WEBHOOKS_URL + webhookType)

		if err != nil {
			return err
		}

		// Verificar si la respuesta fue exitosa (códigos 2xx)
		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			return fmt.Errorf("webhook request failed with status %d: %s", resp.StatusCode(), resp.String())
		}

		return nil
	}
}
