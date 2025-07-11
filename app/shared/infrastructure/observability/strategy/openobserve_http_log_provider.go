package strategy

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transport-app/app/shared/configuration"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func OpenObserveHTTPLogProvider(conf configuration.Conf) (*slog.Logger, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var endpoint OpenObserveHttpEndpoint = OpenObserveHttpEndpoint(os.Getenv(OPENOBSERVE_HTTP_ENDPOINT))

	exporter, err := otlploghttp.New(context.TODO(),
		otlploghttp.WithEndpoint(endpoint.GetDNS()),
		otlploghttp.WithURLPath(endpoint.GetPath()+"/v1/logs"),
		otlploghttp.WithTimeout(5*time.Second),
		otlploghttp.WithHeaders(map[string]string{
			"Authorization": os.Getenv(OPENOBSERVE_AUTHORIZATION),
			"stream-name":   os.Getenv(OPENOBSERVE_STREAM_NAME),
		}),
	)

	if err != nil {
		fmt.Println("Error creating HTTP OTLP log exporter:", err)
		cancel()
		return nil, err
	}

	// Define resource attributes
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(conf.PROJECT_NAME),
		semconv.ServiceVersionKey.String(conf.VERSION),
		semconv.DeploymentEnvironmentKey.String(conf.ENVIRONMENT),
	)

	// Create LoggerProvider
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()
		if err := loggerProvider.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown logger provider:", err)
		}
		cancel()
	}()

	// Create the slog.Logger using the otelslog bridge
	logger := otelslog.NewLogger(
		"openobserve",
		otelslog.WithLoggerProvider(loggerProvider),
	)

	return logger, nil
}
