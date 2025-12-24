package sharedcontext

import (
	"context"

	"micartapro/app/shared/constants"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.opentelemetry.io/otel/trace"
)

type idempotencyKeyContextKey struct{}

var idempotencyKeyKey = &idempotencyKeyContextKey{}

func ContextFromCloudEvent(
	ctx context.Context,
	event cloudevents.Event,
) context.Context {
	traceIDStr, ok1 := event.Extensions()[constants.CloudEventExtensionTraceID].(string)
	spanIDStr, ok2 := event.Extensions()[constants.CloudEventExtensionSpanID].(string)

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

	ctx = trace.ContextWithSpanContext(ctx, sc)

	// Restaurar idempotency key del CloudEvent al contexto
	if idempotencyKeyStr, ok := event.Extensions()[constants.CloudEventExtensionIdempotencyKey].(string); ok && idempotencyKeyStr != "" {
		ctx = context.WithValue(ctx, idempotencyKeyKey, idempotencyKeyStr)
	}

	return ctx
}

// WithIdempotencyKey agrega la idempotency key al contexto.
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKeyKey, key)
}

// IdempotencyKeyFromContext extrae la idempotency key del contexto.
// Retorna el valor de la idempotency key y un bool indicando si existe.
func IdempotencyKeyFromContext(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(idempotencyKeyKey).(string)
	if !ok || key == "" {
		return "", false
	}
	return key, true
}
