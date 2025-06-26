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

type Optimize func(ctx context.Context, request request.OptimizeFleetRequest) (domain.Plan, error)

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
	return func(ctx context.Context, req request.OptimizeFleetRequest) (domain.Plan, error) {
		// Convertir el request a la estructura del dominio de optimización
		fleetOptimization := req.Map()

		vroomRequest, err := mapper.MapOptimizationRequest(ctx, fleetOptimization)
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

		plan := vroomResponse.Map(ctx, fleetOptimization)

		return plan, nil
	}
}
