package subscriber

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.opentelemetry.io/otel/trace"
)

func contextFromCloudEvent(
	ctx context.Context,
	event cloudevents.Event,
) context.Context {
	traceIDStr, ok1 := event.Extensions()["trace_id"].(string)
	spanIDStr, ok2 := event.Extensions()["span_id"].(string)

	if !ok1 || !ok2 {
		return ctx // 👈 silencio total
	}

	traceID, err := trace.TraceIDFromHex(traceIDStr)
	if err != nil {
		return ctx
	}

	spanID, err := trace.SpanIDFromHex(spanIDStr)
	if err != nil {
		return ctx
	}

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
		Remote:  true,
	})

	return trace.ContextWithSpanContext(ctx, sc)
}
