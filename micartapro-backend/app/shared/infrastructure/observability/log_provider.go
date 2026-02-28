package observability

import (
	"log/slog"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability/strategy"

	ioc "github.com/Ignaciojeria/ioc"
)

func init() {
	ioc.Register(newLoggerProvider)
}

func newLoggerProvider(conf configuration.Conf) (*slog.Logger, error) {
	// Get the observability strategy
	observabilityStrategyKey := configuration.Getenv(strategy.OBSERVABILITY_STRATEGY)
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.OpenObserveGRPCLogProvider(conf)
	case "datadog":
		return strategy.DatadogGRPCLogProvider(conf)
	default:
		return strategy.NoOpStdoutLogProvider(conf), nil
	}
}
