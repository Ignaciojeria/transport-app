package observability

import (
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability/strategy"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	otelmeter "go.opentelemetry.io/otel/metric"
)

func init() {
	ioc.Registry(
		newMeterProvider,
		configuration.NewConf,
	)
}

func newMeterProvider(conf configuration.Conf) (otelmeter.Meter, error) {
	// Get the observability strategy
	observabilityStrategyKey := configuration.Getenv(strategy.OBSERVABILITY_STRATEGY)
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.NoOpMeterProvider(conf)
	default:
		return strategy.NoOpMeterProvider(conf)
	}
}
