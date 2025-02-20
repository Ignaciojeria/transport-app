package locationiq

import (
	"context"
	"encoding/json"
	"fmt"

	"transport-app/app/adapter/out/restyclient/locationiq/request"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/httpresty"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type LocationIqOptimization func(
	ctx context.Context,
	plan domain.Plan) (domain.Plan, error)

func init() {
	ioc.Registry(
		NewLocationIqOptimization,
		configuration.NewConf,
		httpresty.NewClient)
}

func NewLocationIqOptimization(conf configuration.Conf, cli *resty.Client) LocationIqOptimization {
	return func(ctx context.Context, plan domain.Plan) (domain.Plan, error) {
		// Create a copy of the original plan to modify
		optimizedPlan := plan
		optimizedPlan.Routes = make([]domain.Route, 0, len(plan.Routes))

		// Map optimization requests
		requests := request.MapOptimizationRequest(conf, plan)

		// Iterate through each route's optimization request
		for _, endpoint := range requests.ENDPOINTS {
			fmt.Println(endpoint)
			// Send optimization request to LocationIQ
			resp, err := cli.R().
				SetContext(ctx).
				Get(endpoint.URL)

			if err != nil {
				return domain.Plan{}, fmt.Errorf("error calling LocationIQ optimization: %w", err)
			}

			// Check response status
			if !resp.IsSuccess() {
				return domain.Plan{}, fmt.Errorf("LocationIQ optimization failed with status %d: %s",
					resp.StatusCode(), string(resp.Body()))
			}

			// Parse the optimization response
			var optimizationResp request.OptimizationResponse
			if err := json.Unmarshal(resp.Body(), &optimizationResp); err != nil {
				return domain.Plan{}, fmt.Errorf("error parsing LocationIQ response: %w", err)
			}

			// Validate response
			if optimizationResp.Code != "Ok" {
				return domain.Plan{}, fmt.Errorf("LocationIQ optimization returned non-OK code: %s", optimizationResp.Code)
			}

			// Map the optimized route
			optimizedRoute := optimizationResp.Map(endpoint.Route)

			// Add the optimized route to the plan
			optimizedPlan.Routes = append(optimizedPlan.Routes, optimizedRoute)
		}

		return optimizedPlan, nil
	}
}
