package vroom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/vroom/mapper"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/hashicorp/go-retryablehttp"
)

type Optimize func(ctx context.Context, request request.OptimizationRequest) (any, error)

func init() {
	ioc.Registry(
		NewOptimize,
		observability.NewObservability,
		NewVroomFastClient,
		NewVroomDefaultClient,
		NewVroomHeavyClient,
		configuration.NewConf,
	)
}
func NewOptimize(
	obs observability.Observability,
	fastClient *retryablehttp.Client,
	defaultClient *retryablehttp.Client,
	heavyClient *retryablehttp.Client,
	conf configuration.Conf,
) Optimize {
	return func(ctx context.Context, request request.OptimizationRequest) (any, error) {

		vroomRequest, err := mapper.MapOptimizationRequest(request)
		if err != nil {
			return nil, err
		}

		// Serializar a JSON
		jsonBytes, err := json.Marshal(vroomRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal VROOM request: %w", err)
		}

		// Log request details
		obs.Logger.InfoContext(ctx,
			"VROOM_REQUEST",
			"url", conf.VROOM_URL,
			"payload", string(jsonBytes),
		)

		// Enviar como io.Reader
		res, err := fastClient.Post(
			conf.VROOM_URL,
			"application/json",
			bytes.NewReader(jsonBytes),
		)
		if err != nil {
			obs.Logger.ErrorContext(ctx,
				"VROOM_REQUEST_ERROR",
				"error", err.Error(),
				"url", conf.VROOM_URL,
			)
			return nil, err
		}

		defer res.Body.Close()

		// Leer el cuerpo de la respuesta
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Log response details
		obs.Logger.InfoContext(ctx,
			"VROOM_RESPONSE",
			"status", res.StatusCode,
			"body", string(bodyBytes),
		)

		// Si el status code no es 200, retornar error con detalles
		if res.StatusCode != 200 {
			obs.Logger.ErrorContext(ctx,
				"VROOM_API_ERROR",
				"status", res.StatusCode,
				"body", string(bodyBytes),
				"request", string(jsonBytes),
			)

			return nil, fmt.Errorf("VROOM API error (status %d): %s\nRequest payload: %s",
				res.StatusCode,
				string(bodyBytes),
				string(jsonBytes))
		}

		return nil, nil
	}
}
