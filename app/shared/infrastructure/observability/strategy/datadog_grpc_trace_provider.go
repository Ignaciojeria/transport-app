package strategy

import (
	"transport-app/app/shared/configuration"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials/insecure"
)

// NewTraceProvider configures a basic trace provider for DataDog.
func DatadogGRPCTraceProvider(conf configuration.Conf) (trace.Tracer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	exporterOpts := []otlptracegrpc.Option{}

	exporterOpts = append(exporterOpts, otlptracegrpc.WithEndpoint(configuration.Getenv(OTEL_EXPORTER_OTLP_ENDPOINT)))
	if configuration.Getenv(OTEL_EXPORTER_OTLP_INSECURE) == "true" {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	}
	client := otlptracegrpc.NewClient(
		exporterOpts...,
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.PROJECT_NAME),
			semconv.DeploymentEnvironmentKey.String(conf.ENVIRONMENT),
		)),
	)
	otel.SetTracerProvider(tp)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Second*2)
		defer shutdownCancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown:", err)
		}
		cancel()
	}()

	return tp.Tracer("observability"), nil
}
