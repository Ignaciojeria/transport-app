package domain

import "context"

type Headers struct {
	Consumer string `json:"consumer"`
	Commerce string `json:"commerce"`
}

func (h Headers) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, h.Commerce, h.Consumer)
}

func (c Headers) UpdateIfChanged(newHeaders Headers) Headers {
	if newHeaders.Consumer != "" {
		c.Consumer = newHeaders.Consumer
	}
	if newHeaders.Commerce != "" {
		c.Commerce = newHeaders.Commerce
	}
	return c
}
