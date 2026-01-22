package sharedcontext

import (
	"context"

	"micartapro/app/shared/constants"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.opentelemetry.io/otel/trace"
)

type idempotencyKeyContextKey struct{}
type userIDContextKey struct{}
type versionIDContextKey struct{}

var idempotencyKeyKey = &idempotencyKeyContextKey{}
var userIDKey = &userIDContextKey{}
var versionIDKey = &versionIDContextKey{}

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

	// Restaurar user ID del CloudEvent al contexto
	if userIDStr, ok := event.Extensions()[constants.CloudEventExtensionUserID].(string); ok && userIDStr != "" {
		ctx = context.WithValue(ctx, userIDKey, userIDStr)
	}

	// Restaurar version ID del CloudEvent al contexto
	if versionIDStr, ok := event.Extensions()[constants.CloudEventExtensionVersionID].(string); ok && versionIDStr != "" {
		ctx = context.WithValue(ctx, versionIDKey, versionIDStr)
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

// WithUserID agrega el user ID al contexto.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extrae el user ID del contexto.
// Retorna el valor del user ID y un bool indicando si existe.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}

// WithVersionID agrega el version ID al contexto.
func WithVersionID(ctx context.Context, versionID string) context.Context {
	return context.WithValue(ctx, versionIDKey, versionID)
}

// VersionIDFromContext extrae el version ID del contexto.
// Retorna el valor del version ID y un bool indicando si existe.
func VersionIDFromContext(ctx context.Context) (string, bool) {
	versionID, ok := ctx.Value(versionIDKey).(string)
	if !ok || versionID == "" {
		return "", false
	}
	return versionID, true
}
