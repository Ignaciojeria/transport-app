package domain

import (
	"context"
)

type Route struct {
	ReferenceID string
	Origin      NodeInfo
	Destination NodeInfo
	Vehicle     Vehicle
	Orders      []Order
	TimeWindow  TimeWindow
	Geometry    RouteGeometry
}

type RouteGeometry struct {
	Encoding string `json:"encoding"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

type TimeWindow struct {
	Start string
	End   string
}

func (s Route) DocID(ctx context.Context) DocumentID {
	return HashByCountry(ctx, s.ReferenceID)
}
