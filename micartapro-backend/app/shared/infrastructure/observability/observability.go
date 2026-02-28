package observability

import (
	"log/slog"

	ioc "github.com/Ignaciojeria/ioc"
	otelmeter "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Observability struct {
	Tracer trace.Tracer
	Logger *slog.Logger
	Meter  otelmeter.Meter
}

func init() {
	ioc.Register(NewObservability)
}
func NewObservability(
	tracer trace.Tracer,
	logger *slog.Logger,
	meter otelmeter.Meter) Observability {
	return Observability{
		Tracer: tracer,
		Logger: logger,
		Meter:  meter,
	}
}
