package observability

import (
	"log/slog"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability/strategy"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newLoggerProvider,
		configuration.NewConf,
	)
}

func newLoggerProvider(conf configuration.Conf) (*slog.Logger, error) {
	// Get the observability strategy
	observabilityStrategyKey := configuration.Getenv(strategy.OBSERVABILITY_STRATEGY)
	switch observabilityStrategyKey {
	case "openobserve":
		return strategy.OpenObserveHTTPLogProvider(conf)
	case "datadog":
		return strategy.DatadogGRPCLogProvider(conf)
	default:
		return strategy.NoOpStdoutLogProvider(conf), nil
	}
}
