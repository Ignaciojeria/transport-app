package fuegoapiclient

import (
	"context"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/httpresty"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type PostUpsertRoute func(ctx context.Context, e any) error

func init() {
	ioc.Registry(
		NewPostUpsertRoute,
		httpresty.NewClient,
		configuration.NewConf,
	)
}

// TODO ver la manera de usar una entidad de dominio en vez de un request
func NewPostUpsertRoute(c *resty.Client, conf configuration.Conf) PostUpsertRoute {
	return func(ctx context.Context, e any) error {
		host := conf.MASTER_NODE_HOST
		apiKey := conf.MASTER_NODE_API_KEY

		url := fmt.Sprintf("%s/routes", host)

		//TODO CONVERTIR ENTIDAD DE DOMINIO A REQUEST
		var req request.UpsertRouteRequest

		res, err := c.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
			SetBody(req).
			Post(url)

		if err != nil {
			return fmt.Errorf("failed to make HTTP request: %w", err)
		}

		if res.IsError() {
			return fmt.Errorf("API error (status %d): %s", res.StatusCode(), res.String())
		}

		return nil
	}
}
