package domain

import "context"

type PlanningStatus struct {
	Value string
}

func (ps PlanningStatus) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, ps.Value)
}
