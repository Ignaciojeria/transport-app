package strategy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transport-app/app/shared/configuration"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func OpenObserveHTTPTraceProvider(conf configuration.Conf) (trace.Tracer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	var endpoint OpenObserveHttpEndpoint = OpenObserveHttpEndpoint(os.Getenv(OPENOBSERVE_HTTP_ENDPOINT))

	otlpHTTPExporter, err := otlptracehttp.New(context.TODO(),
		otlptracehttp.WithEndpoint(endpoint.GetDNS()),
		otlptracehttp.WithURLPath(endpoint.GetPath()+"/v1/traces"),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": os.Getenv(OPENOBSERVE_AUTHORIZATION),
			"stream-name":   os.Getenv(OPENOBSERVE_STREAM_NAME),
		}),
	)
	if err != nil {
		fmt.Println("Error creating HTTP OTLP exporter:", err)
		cancel()
		return nil, nil
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(conf.PROJECT_NAME),
		semconv.ServiceVersionKey.String(conf.VERSION),
		attribute.String("environment", conf.ENVIRONMENT),
	)

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpHTTPExporter),
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
