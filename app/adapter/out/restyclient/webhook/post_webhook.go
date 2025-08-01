package webhook

import (
	"context"
	"fmt"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpresty"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type PostWebhook func(ctx context.Context, input domain.Webhook) error

func init() {
	ioc.Registry(NewPostWebhook, httpresty.NewClient)
}

func NewPostWebhook(c *resty.Client) PostWebhook {
	return func(ctx context.Context, webhook domain.Webhook) error {
		// Validar el webhook antes de enviarlo
		if err := webhook.Validate(); err != nil {
			return err
		}

		// Crear la petición HTTP con resty
		resp, err := c.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeaders(webhook.Headers).
			SetBody(webhook.Body).
			Post(webhook.URL)

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
