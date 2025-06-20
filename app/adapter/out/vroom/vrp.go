package vroom

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/vroom/mapper"
	"transport-app/app/adapter/out/vroom/model"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-resty/resty/v2"
)

type Optimize func(ctx context.Context, request request.FleetsOptimizationRequest) (domain.Plan, error)

func init() {
	ioc.Registry(
		NewOptimize,
		observability.NewObservability,
		NewVroomRestyFastClient,
		configuration.NewConf,
	)
}

func NewOptimize(
	obs observability.Observability,
	restyClient *resty.Client,
	conf configuration.Conf,
) Optimize {
	return func(ctx context.Context, req request.FleetsOptimizationRequest) (domain.Plan, error) {
		vroomRequest, err := mapper.MapOptimizationRequest(ctx, req)
		if err != nil {
			return domain.Plan{}, err
		}

		jsonBytes, err := json.Marshal(vroomRequest)
		if err != nil {
			return domain.Plan{}, fmt.Errorf("failed to marshal VROOM request: %w", err)
		}

		obs.Logger.InfoContext(ctx,
			"VROOM_REQUEST",
			"url", conf.VROOM_URL,
			"payload", string(jsonBytes),
		)

		res, err := restyClient.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetBody(jsonBytes).
			Post(conf.VROOM_URL)

		if err != nil {
			obs.Logger.ErrorContext(ctx,
				"VROOM_REQUEST_ERROR",
				"error", err.Error(),
				"url", conf.VROOM_URL,
			)
			return domain.Plan{}, err
		}

		obs.Logger.InfoContext(ctx,
			"VROOM_RESPONSE",
			"status", res.StatusCode(),
			"body", res.String(),
		)

		if res.IsError() {
			obs.Logger.ErrorContext(ctx,
				"VROOM_API_ERROR",
				"status", res.StatusCode(),
				"body", res.String(),
				"request", string(jsonBytes),
			)

			return domain.Plan{}, fmt.Errorf("VROOM API error (status %d): %s\nRequest payload: %s",
				res.StatusCode(),
				res.String(),
				string(jsonBytes))
		}

		// Deserializar la respuesta de VROOM
		var vroomResponse model.VroomOptimizationResponse
		if err := json.Unmarshal(res.Body(), &vroomResponse); err != nil {
			obs.Logger.ErrorContext(ctx,
				"VROOM_RESPONSE_DESERIALIZATION_ERROR",
				"error", err.Error(),
				"body", res.String(),
			)
			return domain.Plan{}, fmt.Errorf("failed to deserialize VROOM response: %w", err)
		}

		// Aplicar el mapper para convertir la respuesta a domain.Plan
		plan := vroomResponse.Map(ctx, req)

		// Log del plan completo para debugging
		planJSON, err := json.MarshalIndent(plan, "", "  ")
		if err != nil {
			obs.Logger.ErrorContext(ctx,
				"VROOM_PLAN_JSON_MARSHAL_ERROR",
				"error", err.Error(),
			)
		} else {
			obs.Logger.InfoContext(ctx,
				"VROOM_PLAN_COMPLETE",
				"plan_json", string(planJSON),
			)
		}

		obs.Logger.InfoContext(ctx,
			"VROOM_OPTIMIZATION_COMPLETED",
			"plan_reference_id", plan.ReferenceID,
			"routes_count", len(plan.Routes),
			"unassigned_orders_count", len(plan.UnassignedOrders),
		)

		return plan, nil
	}
}
