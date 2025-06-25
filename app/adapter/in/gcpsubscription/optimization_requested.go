package gcpsubscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/gcpsubscription/mapper"
	"transport-app/app/adapter/out/vroom"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newOptimizationRequested,
		subscriptionwrapper.NewSubscriptionManager,
		configuration.NewConf,
		vroom.NewOptimize,
		observability.NewObservability,
		usecase.NewCreatePlan,
	)
}
func newOptimizationRequested(
	sm subscriptionwrapper.SubscriptionManager,
	conf configuration.Conf,
	optimize vroom.Optimize,
	observability observability.Observability,
	createPlan usecase.CreatePlan,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := conf.OPTIMIZATION_REQUESTED_SUBSCRIPTION
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 15
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		if m.Attributes["eventType"] != "optimizationRequested" {
			m.Ack()
			return http.StatusAccepted, nil
		}
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m.Attributes))
		var input request.OptimizeFleetRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		res, err := optimize(ctx, input)

		if err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		plan := mapper.MapOptimizedFleetToPlan(ctx, res)

		if err := createPlan(ctx, plan); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		// Obtener el FleetOptimization original para las horas de inicio de los vehículos
		originalFleet := input.Map()
		vroom.ExportPolylineJSONFromOptimizedFleet("ui/static/dev/polyline.json", res, originalFleet)

		observability.Logger.InfoContext(ctx, "Optimization requested", "res", res)
		fmt.Println("works")
		m.Ack()
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		Start(subscriptionRef)
	return messageProcessor
}
