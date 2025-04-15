package domain

import "context"

type Headers struct {
	Consumer string `json:"consumer"`
	Commerce string `json:"commerce"`
}

func (h Headers) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, h.Commerce, h.Consumer)
}
