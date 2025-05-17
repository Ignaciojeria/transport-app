package domain

import "context"

type Route struct {
	ReferenceID string
	Origin      NodeInfo
	Destination NodeInfo
	Vehicle     Vehicle
	Orders      []Order
}

func (s Route) DocID(ctx context.Context) DocumentID {
	return HashByCountry(ctx, s.ReferenceID)
}
