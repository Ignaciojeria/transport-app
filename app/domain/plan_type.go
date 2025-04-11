package domain

import "context"

type PlanType struct {
	Value string
}

func (pt PlanType) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, pt.Value)
}
