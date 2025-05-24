package domain

import "context"

type SizeCategory struct {
	Code string `json:"code"`
}

func (s SizeCategory) DocumentID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, s.Code)
}
