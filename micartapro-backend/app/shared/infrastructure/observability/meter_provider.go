package observability

import (
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability/strategy"

	ioc "github.com/Ignaciojeria/ioc"
	otelmeter "go.opentelemetry.io/otel/metric"
)

func init() {
	ioc.Register(newMeterProvider)
}

func newMeterProvider(conf configuration.Conf) (otelmeter.Meter, error) {
	// Get the observability strategy
	observabilityStrategyKey := configuration.Getenv(strategy.OBSERVABILITY_STRATEGY)
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.NewGRPCOpenObserveMeterProvider(conf)
	default:
		return strategy.NoOpMeterProvider(conf)
	}
}
