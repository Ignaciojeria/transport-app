package domain

import "context"

type Route struct {
	ReferenceID     string
	PlanReferenceID string
	TimeWindow      TimeWindow
	Origin          NodeInfo
	Destination     NodeInfo
	Vehicle         Vehicle
	Orders          []Order
}

// TimeWindow representa una ventana de tiempo
type TimeWindow struct {
	Start string
	End   string
}

func (s Route) DocID(ctx context.Context) DocumentID {
	return HashByCountry(ctx, s.ReferenceID)
}
