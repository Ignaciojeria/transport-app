package observability

import (
	"log/slog"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel"
	otelmeter "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Observability struct {
	Tracer trace.Tracer
	Logger *slog.Logger
	Meter  otelmeter.Meter
}

func init() {
	ioc.Registry(
		NewObservability,
		newTraceProvider,
		newLoggerProvider,
		newMeterProvider)
}

func NewObservability(
	tracer trace.Tracer,
	logger *slog.Logger,
	meter otelmeter.Meter) Observability {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return Observability{
		Tracer: tracer,
		Logger: logger,
		Meter:  meter,
	}
}
