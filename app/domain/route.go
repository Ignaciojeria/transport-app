package domain

import "context"

type Route struct {
	ReferenceID string
	Destination NodeInfo
	Vehicle     Vehicle
	Operator    Operator
	Orders      []Order
}

func (s Route) DocID(ctx context.Context) DocumentID {
	return HashByCountry(ctx, s.ReferenceID)
}
