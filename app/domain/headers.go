package domain

import (
	"context"
	"transport-app/app/shared/sharedcontext"

	"go.opentelemetry.io/otel/baggage"
)

type Headers struct {
	Consumer string `json:"consumer"`
	Commerce string `json:"commerce"`
	Channel  string `json:"channel"`
}

func (h Headers) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, h.Commerce, h.Consumer, h.Channel)
}

func (h *Headers) SetFromContext(ctx context.Context) {
	bag := baggage.FromContext(ctx)
	commerce := bag.Member(sharedcontext.BaggageCommerce).Value()
	h.Commerce = commerce
	consumer := bag.Member(sharedcontext.BaggageConsumer).Value()
	h.Consumer = consumer
	channel := bag.Member(sharedcontext.BaggageChannel).Value()
	h.Channel = channel
}

func (h Headers) UpdateIfChanged(in Headers) (Headers, bool) {
	changed := false

	if in.Commerce != "" && in.Commerce != h.Commerce {
		h.Commerce = in.Commerce
		changed = true
	}

	if in.Consumer != "" && in.Consumer != h.Consumer {
		h.Consumer = in.Consumer
		changed = true
	}

	if in.Channel != "" && in.Channel != h.Channel {
		h.Channel = in.Channel
		changed = true
	}

	return h, changed
}

func (h Headers) IsEmpty() bool {
	return h.Commerce == "" && h.Consumer == "" && h.Channel == ""
}
