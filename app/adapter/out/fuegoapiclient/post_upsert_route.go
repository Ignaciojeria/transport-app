package fuegoapiclient

import (
	"context"
	"fmt"

	"transport-app/app/adapter/out/fuegoapiclient/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/httpresty"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type PostUpsertRoute func(ctx context.Context, e domain.Route) error

func init() {
	ioc.Registry(
		NewPostUpsertRoute,
		httpresty.NewClient,
		configuration.NewConf,
	)
}

// TODO ver la manera de usar una entidad de dominio en vez de un request
func NewPostUpsertRoute(c *resty.Client, conf configuration.Conf) PostUpsertRoute {
	return func(ctx context.Context, e domain.Route) error {
		host := conf.MASTER_NODE_HOST
		apiKey := conf.MASTER_NODE_API_KEY

		url := fmt.Sprintf("%s/routes", host)

		req := mapper.MapUpsertRouteRequest(e)

		res, err := c.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("channel", sharedcontext.ChannelFromContext(ctx)).
			SetHeader("tenant", sharedcontext.TenantIDFromContext(ctx).String()).
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
