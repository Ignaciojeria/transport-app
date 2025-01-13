package observability

import (
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability/strategy"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	ioc.Registry(
		newTraceProvider,
		configuration.NewConf,
	)
}

// RegisterTraceProvider determines whether to use OpenObserve, Datadog or non provider based on the existing environment variables.
func newTraceProvider(conf configuration.Conf) (trace.Tracer, error) {
	// Get the observability strategy
	observabilityStrategyKey := configuration.Getenv(strategy.OBSERVABILITY_STRATEGY)
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.OpenObserveHTTPTraceProvider(conf)
	case "datadog":
		return strategy.DatadogGRPCTraceProvider(conf)
	default:
		return strategy.NoOpTraceProvider(conf)
	}
}
