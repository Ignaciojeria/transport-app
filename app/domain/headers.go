package domain

import (
	"context"
	"transport-app/app/shared/sharedcontext"

	"go.opentelemetry.io/otel/baggage"
)

type Headers struct {
	Consumer string `json:"consumer"`
	Commerce string `json:"commerce"`
}

func (h Headers) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, h.Commerce, h.Consumer)
}

func (h *Headers) SetFromContext(ctx context.Context) {
	bag := baggage.FromContext(ctx)
	commerce := bag.Member(sharedcontext.BaggageCommerce).Value()
	h.Commerce = commerce
	consumer := bag.Member(sharedcontext.BaggageConsumer).Value()
	h.Consumer = consumer
}
